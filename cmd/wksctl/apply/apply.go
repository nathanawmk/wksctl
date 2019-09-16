package apply

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/weaveworks/wksctl/pkg/addons"
	"github.com/weaveworks/wksctl/pkg/apis/wksprovider/machine/config"
	wksos "github.com/weaveworks/wksctl/pkg/apis/wksprovider/machine/os"
	"github.com/weaveworks/wksctl/pkg/manifests"
	"github.com/weaveworks/wksctl/pkg/specs"
	"github.com/weaveworks/wksctl/pkg/utilities/kubeadm"
	"github.com/weaveworks/wksctl/pkg/utilities/manifest"
	"github.com/weaveworks/wksctl/pkg/version"
)

// Cmd represents the apply command
var Cmd = &cobra.Command{
	Use:   "apply",
	Short: "Create or update a Kubernetes cluster",
	RunE:  func(_ *cobra.Command, _ []string) error { a := Applier{&globalParams}; return a.Apply() },
}

type Params struct {
	clusterManifestPath  string
	machinesManifestPath string
	controllerImage      string
	gitURL               string
	gitBranch            string
	gitPath              string
	gitDeployKeyPath     string
	sealedSecretKeyPath  string
	sealedSecretCertPath string
	configDirectory      string
	namespace            string
	useManifestNamespace bool
	verbose              bool
}

var globalParams Params

func init() {
	Cmd.Flags().StringVar(&globalParams.clusterManifestPath, "cluster", "cluster.yaml", "Location of cluster manifest")
	Cmd.Flags().StringVar(&globalParams.machinesManifestPath, "machines", "machines.yaml", "Location of machines manifest")
	Cmd.Flags().StringVar(&globalParams.gitURL, "git-url", "", "Git repo containing your cluster and machine information")
	Cmd.Flags().StringVar(&globalParams.gitBranch, "git-branch", "master", "Git branch WKS should use to sync with your cluster")
	Cmd.Flags().StringVar(&globalParams.gitPath, "git-path", ".", "Relative path to files in Git")
	Cmd.Flags().StringVar(&globalParams.gitDeployKeyPath, "git-deploy-key", "", "Path to the Git deploy key")
	Cmd.Flags().StringVar(&globalParams.sealedSecretKeyPath, "sealed-secret-key", "", "Path to a key used to decrypt sealed secrets")
	Cmd.Flags().StringVar(&globalParams.sealedSecretCertPath, "sealed-secret-cert", "", "Path to a certificate used to encrypt sealed secrets")
	Cmd.Flags().StringVar(&globalParams.configDirectory, "config-directory", ".", "Directory containing configuration information for the cluster")
	Cmd.Flags().StringVar(&globalParams.namespace, "namespace", manifest.DefaultNamespace, "namespace override for WKS components")
	Cmd.Flags().BoolVar(&globalParams.useManifestNamespace, "use-manifest-namespace", false, "use namespaces from supplied manifests (overriding any --namespace argument)")

	// Intentionally shadows the globally defined --verbose flag.
	Cmd.Flags().BoolVar(&globalParams.verbose, "verbose", false, "Enable verbose output")

	// Hide controller-image flag as it is a helper/debug flag.
	Cmd.Flags().StringVar(&globalParams.controllerImage, "controller-image", "quay.io/wksctl/controller:"+version.ImageTag, "Controller image override")
	Cmd.Flags().MarkHidden("controller-image")
}

type Applier struct {
	Params *Params
}

func (a *Applier) Apply() error {
	// Default to using the git deploy key to decrypt sealed secrets
	if a.Params.sealedSecretKeyPath == "" && a.Params.gitDeployKeyPath != "" {
		a.Params.sealedSecretKeyPath = a.Params.gitDeployKeyPath
	}

	var closer func()
	var err error
	cpath := filepath.Join(a.Params.gitPath, a.Params.clusterManifestPath)
	mpath := filepath.Join(a.Params.gitPath, a.Params.machinesManifestPath)

	cpath, mpath, closer, err = manifests.Get(cpath, mpath, a.Params.gitURL, a.Params.gitBranch, a.Params.gitDeployKeyPath, a.Params.gitPath)
	if err != nil {
		return err
	}
	return a.initiateCluster(cpath, mpath, closer)
}

func (a *Applier) initiateCluster(clusterManifestPath, machinesManifestPath string, closer func()) error {
	defer closer()
	sp := specs.NewFromPaths(clusterManifestPath, machinesManifestPath)
	sshClient, err := sp.GetSSHClient(a.Params.verbose)
	if err != nil {
		return errors.Wrap(err, "failed to create SSH client")
	}
	defer sshClient.Close()
	installer, err := wksos.Identify(sshClient)
	if err != nil {
		return errors.Wrapf(err, "failed to identify operating system for seed node (%s)", sp.GetMasterPublicAddress())
	}

	// N.B.: we generate this bootstrap token where wksctl apply is run hoping
	// that this will be on a machine which has been running for a while, and
	// therefore will generate a "more random" token, than we would on a
	// potentially newly created VM which doesn't have much entropy yet.
	token, err := kubeadm.GenerateBootstrapToken()
	if err != nil {
		return errors.Wrap(err, "failed to generate bootstrap token")
	}

	// Point config dir at sync repo if using github and the user didn't override it
	configDir := a.Params.configDirectory
	if configDir == "." && a.Params.gitURL != "" {
		configDir = filepath.Dir(clusterManifestPath)
	}

	ns := ""
	if !a.Params.useManifestNamespace {
		ns = a.Params.namespace
	}

	// TODO(damien): Transform the controller image into an addon.
	controllerImage, err := addons.UpdateImage(a.Params.controllerImage, sp.ClusterSpec.ImageRepository)
	if err != nil {
		errors.Wrap(err, "failed to apply the cluster's image repository to the WKS controller's image")
	}
	if err := installer.SetupSeedNode(wksos.SeedNodeParams{
		PublicIP:             sp.GetMasterPublicAddress(),
		PrivateIP:            sp.GetMasterPrivateAddress(),
		ClusterManifestPath:  clusterManifestPath,
		MachinesManifestPath: machinesManifestPath,
		SSHKeyPath:           sp.GetSSHKeyPath(),
		BootstrapToken:       token,
		KubeletConfig: config.KubeletConfig{
			NodeIP:        sp.GetMasterPrivateAddress(),
			CloudProvider: sp.GetCloudProvider(),
		},
		ControllerImageOverride: controllerImage,
		GitData: wksos.GitParams{
			GitURL:           a.Params.gitURL,
			GitBranch:        a.Params.gitBranch,
			GitPath:          a.Params.gitPath,
			GitDeployKeyPath: a.Params.gitDeployKeyPath,
		},
		SealedSecretKeyPath:  a.Params.sealedSecretKeyPath,
		SealedSecretCertPath: a.Params.sealedSecretCertPath,
		ConfigDirectory:      configDir,
		ImageRepository:      sp.ClusterSpec.ImageRepository,
		ExternalLoadBalancer: sp.ClusterSpec.APIServer.ExternalLoadBalancer,
		AdditionalSANs:       sp.ClusterSpec.APIServer.AdditionalSANs,
		Namespace:            ns,
	}); err != nil {
		return errors.Wrapf(err, "failed to set up seed node (%s)", sp.GetMasterPublicAddress())
	}

	return nil
}