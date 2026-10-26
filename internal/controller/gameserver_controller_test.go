package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	respawnv1alpha1 "github.com/bblaker/respawn/api/v1alpha1"
)

var _ = Describe("GameServer Controller", func() {
	Context("When reconciling a resource", func() {
		const (
			resourceName      = "test-resource"
			resourceNamespace = "default"
		)

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: resourceNamespace,
		}
		gameserver := &respawnv1alpha1.GameServer{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind GameServer")
			err := k8sClient.Get(ctx, typeNamespacedName, gameserver)
			if err != nil && errors.IsNotFound(err) {
				resource := newGameServer(resourceName)
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &respawnv1alpha1.GameServer{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance GameServer")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &GameServerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
