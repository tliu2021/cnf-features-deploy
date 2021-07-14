package utils

import (
	"os"

	sriovNamespaces "github.com/k8snetworkplumbingwg/sriov-network-operator/test/util/namespaces"
	perfUtils "github.com/openshift-kni/performance-addon-operators/functests/utils"
	ptpUtils "github.com/openshift/ptp-operator/test/utils"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift-kni/cnf-features-deploy/cnf-tests/testsuites/pkg/k8sreporter"
	"github.com/openshift-kni/cnf-features-deploy/cnf-tests/testsuites/pkg/namespaces"

	gkopv1alpha "github.com/gatekeeper/gatekeeper-operator/api/v1alpha1"
	sriovv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
	n3000v1 "github.com/open-ness/openshift-operator/N3000/api/v1"
	sriovfecv1 "github.com/open-ness/openshift-operator/sriov-fec/api/v1"
	gkv1alpha "github.com/open-policy-agent/gatekeeper/apis/mutations/v1alpha1"
	performancev2 "github.com/openshift-kni/performance-addon-operators/api/v2"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	ptpv1 "github.com/openshift/ptp-operator/api/v1"
)

// NewReporter creates a specific reporter for CNF tests
func NewReporter(reportPath string) (*k8sreporter.KubernetesReporter, error) {
	addToScheme := func(s *runtime.Scheme) {
		ptpv1.AddToScheme(s)
		mcfgv1.AddToScheme(s)
		performancev2.SchemeBuilder.AddToScheme(s)
		sriovv1.AddToScheme(s)
		gkv1alpha.AddToScheme(s)
		gkopv1alpha.AddToScheme(s)
	}

	namespacesToDump := map[string]string{
		namespaces.PerformanceOperator:      "performance",
		namespaces.PTPOperator:              "ptp",
		namespaces.SRIOVOperator:            "sriov",
		NamespaceTesting:                    "other",
		perfUtils.NamespaceTesting:          "performance",
		namespaces.DpdkTest:                 "dpdk",
		sriovNamespaces.Test:                "sriov",
		ptpUtils.NamespaceTesting:           "ptp",
		namespaces.SCTPTest:                 "sctp",
		namespaces.XTU32Test:                "xt_u32",
		namespaces.IntelOperator:            "intel",
		namespaces.OVSQOSTest:               "ovs_qos",
		GatekeeperNamespace:                 "gatekeeper",
		OperatorNamespace:                   "gatekeeper",
		GatekeeperTestingNamespace:          "gatekeeper",
		GatekeeperMutationIncludedNamespace: "gatekeeper",
		GatekeeperMutationExcludedNamespace: "gatekeeper",
		GatekeeperMutationEnabledNamespace:  "gatekeeper",
		GatekeeperMutationDisabledNamespace: "gatekeeper",
		GatekeeperTestObjectNamespace:       "gatekeeper",
	}

	crds := []k8sreporter.CRData{
		{Cr: &mcfgv1.MachineConfigPoolList{}},
		{Cr: &mcfgv1.MachineConfigList{}},
		{Cr: &ptpv1.PtpConfigList{}},
		{Cr: &ptpv1.NodePtpDeviceList{}},
		{Cr: &ptpv1.PtpOperatorConfigList{}},
		{Cr: &performancev2.PerformanceProfileList{}},
		{Cr: &sriovv1.SriovNetworkNodePolicyList{}},
		{Cr: &sriovv1.SriovNetworkList{}},
		{Cr: &sriovv1.SriovNetworkNodeStateList{}},
		{Cr: &sriovv1.SriovOperatorConfigList{}},
		{Cr: &sriovfecv1.SriovFecNodeConfigList{}},
		{Cr: &sriovfecv1.SriovFecClusterConfigList{}},
		{Cr: &n3000v1.N3000ClusterList{}},
		{Cr: &n3000v1.N3000NodeList{}},
		{Cr: &gkv1alpha.AssignList{}},
		{Cr: &gkv1alpha.AssignMetadataList{}},
		{Cr: &gkopv1alpha.GatekeeperList{}},
	}

	skipByNamespace := func(ns string) bool {
		_, found := namespacesToDump[ns]
		return !found
	}

	err := os.Mkdir(reportPath, 0755)
	if err != nil {
		return nil, err
	}

	res, err := k8sreporter.New("", addToScheme, skipByNamespace, reportPath, crds...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
