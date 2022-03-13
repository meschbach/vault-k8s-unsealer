package main

import (
	"context"
	"encoding/json"
	"fmt"
	vault "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

type runOptions struct {
	inClusterConfig bool
	namespace       string
	podName         string
	secretName      string
}

type VaultSealStore struct {
	Keys      []string `json:"keys"`
	RootToken string   `json:"root_token"`
}

func watchVaultPod(ctx context.Context, clientset *kubernetes.Clientset, opts runOptions) error {
	//TODO: optimize what is being watched and manage cache
	watch, err := clientset.CoreV1().Pods(opts.namespace).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	fmt.Printf("Watching for changes in %s/%s\n", opts.namespace, opts.podName)
	for _ = range watch.ResultChan() {
		fmt.Printf("Change detected.  Checking pod status.\n")
		pod, err := clientset.CoreV1().Pods(opts.namespace).Get(ctx, opts.podName, metav1.GetOptions{})
		if err != nil {
			watch.Stop()
			return err
		}

		fmt.Printf("Pod IP %s\n", pod.Status.PodIP)
		config := vault.DefaultConfig()
		config.Address = "http://" + pod.Status.PodIP + ":8200"
		client, err := vault.NewClient(config)
		if err != nil {
			return err
		}
		sys := client.Sys()
		status, err := sys.SealStatus()
		if err != nil {
			fmt.Printf("Failed to connect with pod: %s\n", err.Error())
			continue
		}

		if !status.Sealed {
			fmt.Printf("Not sealed, listening for further changes\n")
			continue
		}

		fmt.Printf("Sealed.  Initialized? %t\n", status.Initialized)
		secret, err := clientset.CoreV1().Secrets(opts.namespace).Get(ctx, opts.secretName, metav1.GetOptions{})
		if err != nil {
			if e, ok := err.(*errors.StatusError); ok {
				if e.Status().Reason == "NotFound" {
					fmt.Printf("Secret not found, awaiting creation.\n")
					continue
				}
			}
			return err
		}

		jsonBytes, ok := secret.Data["json"]
		if !ok {
			fmt.Printf("missing key json in secret\n")
			continue
		}
		var sealInfo VaultSealStore
		if err := json.Unmarshal(jsonBytes, &sealInfo); err != nil {
			fmt.Printf("error parsing: %e\n", err)
			continue
		}

		for _, k := range sealInfo.Keys {
			_, err := sys.Unseal(k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	var runOpts runOptions
	inCluster := &cobra.Command{
		Use:   "in-cluster",
		Short: "run using a mounted secret for unsealing",
		RunE: func(cmd *cobra.Command, args []string) error {
			// creates the in-cluster config
			config, err := rest.InClusterConfig()
			if err != nil {
				return err
			}
			// creates the clientset
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				return err
			}
			return watchVaultPod(cmd.Context(), clientset, runOpts)
		},
	}

	outCluster := &cobra.Command{
		Use:   "out-cluster",
		Short: "run using a mounted secret for unsealing",
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			kubeConfigPath, maybe := os.LookupEnv("KUBECONFIG")
			if !maybe {
				var homeDir string
				if homeDir, err = os.UserHomeDir(); err != nil {
					return err
				}
				kubeConfigPath = filepath.Join(homeDir, ".kube", "config")
			}
			kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
			if err != nil {
				return err
			}
			clientset, err := kubernetes.NewForConfig(kubeConfig)
			if err != nil {
				return err
			}
			return watchVaultPod(cmd.Context(), clientset, runOpts)
		},
	}

	root := &cobra.Command{
		Use:          "k8s-watcher",
		Short:        "Watches a given k8s pod name for changes then queries to ensure unsealed",
		SilenceUsage: true,
	}
	root.PersistentFlags().StringVar(&runOpts.namespace, "namespace", os.Getenv("VAULT_NAMESPACE"), "target vault namespace to observe")
	root.PersistentFlags().StringVar(&runOpts.podName, "pod-name", os.Getenv("VAULT_POD"), "target vault pod to observe")
	root.PersistentFlags().StringVar(&runOpts.secretName, "secret-name", os.Getenv("VAULT_SECRET"), "secret name to store and retrieve unseal keys")
	root.AddCommand(inCluster)
	root.AddCommand(outCluster)

	if err := root.Execute(); err != nil {
		os.Exit(-1)
	}
}
