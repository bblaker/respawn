package controller

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"

	respawnv1alpha1 "github.com/bblaker/respawn/api/v1alpha1"
)

// These specs exercise the generated CRD schemas against a real API server:
// defaults are applied and basic OpenAPI validation rejects bad input.
// Cross-field rules (CEL) arrive in M2 with their own specs.
var _ = Describe("CRD schemas", func() {
	It("applies GameServer defaults", func() {
		gs := newGameServer("schema-defaults")
		Expect(k8sClient.Create(ctx, gs)).To(Succeed())
		DeferCleanup(func() { Expect(k8sClient.Delete(ctx, gs)).To(Succeed()) })

		got := &respawnv1alpha1.GameServer{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: gs.Name, Namespace: testNamespace}, got)).To(Succeed())

		Expect(got.Spec.Minecraft.Type).To(Equal(respawnv1alpha1.MinecraftTypeVanilla))
		Expect(got.Spec.Minecraft.Memory.Equal(resource.MustParse("2Gi"))).To(BeTrue())
		Expect(got.Spec.Minecraft.MaxPlayers).To(Equal(int32(20)))
		Expect(got.Spec.Storage.ReclaimPolicy).To(Equal(respawnv1alpha1.ReclaimRetain))
		Expect(got.Spec.Power.Mode).To(Equal(respawnv1alpha1.PowerAlwaysOn))
		Expect(got.Spec.Power.IdleTimeout.Duration).To(Equal(15 * time.Minute))
		Expect(got.Spec.Power.PollInterval.Duration).To(Equal(30 * time.Second))
		Expect(got.Spec.Power.StartTimeout.Duration).To(Equal(10 * time.Minute))
		Expect(got.Spec.Power.WakeCooldown.Duration).To(Equal(time.Minute))
		Expect(string(got.Spec.Service.Type)).To(Equal("LoadBalancer"))
	})

	It("applies backup defaults", func() {
		gs := newGameServer("schema-backup-defaults")
		gs.Spec.Backup = &respawnv1alpha1.BackupSpec{}
		gs.Spec.Backup.RepositoryRef.Name = "repo"
		Expect(k8sClient.Create(ctx, gs)).To(Succeed())
		DeferCleanup(func() { Expect(k8sClient.Delete(ctx, gs)).To(Succeed()) })

		Expect(gs.Spec.Backup.OnDelete).To(Equal(new(true)))
		Expect(gs.Spec.Backup.MinInterval.Duration).To(Equal(time.Hour))
	})

	It("applies GameServerBackup defaults", func() {
		b := newGameServerBackup("schema-backup", "some-server")
		Expect(k8sClient.Create(ctx, b)).To(Succeed())
		DeferCleanup(func() { Expect(k8sClient.Delete(ctx, b)).To(Succeed()) })

		Expect(b.Spec.DeletionPolicy).To(Equal(respawnv1alpha1.BackupRetain))
	})

	DescribeTable("rejects invalid GameServers",
		func(mutate func(*respawnv1alpha1.GameServer)) {
			gs := newGameServer("schema-invalid")
			mutate(gs)
			err := k8sClient.Create(ctx, gs)
			Expect(apierrors.IsInvalid(err)).To(BeTrue(), "expected Invalid, got %v", err)
		},
		Entry("unknown power mode", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Power.Mode = "Sometimes"
		}),
		Entry("unknown Minecraft type", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Minecraft.Type = "BUKKIT"
		}),
		Entry("empty Minecraft version", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Minecraft.Version = ""
		}),
		Entry("service port out of range", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Service.Port = new(int32(70000))
		}),
		Entry("unknown reclaim policy", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Storage.ReclaimPolicy = "Recycle"
		}),
		Entry("malformed snapshot ID", func(gs *respawnv1alpha1.GameServer) {
			gs.Spec.Storage.RestoreFrom = &respawnv1alpha1.RestoreSource{SnapshotID: "not-hex!"}
		}),
	)

	It("rejects a BackupRepository endpoint without a scheme", func() {
		repo := newBackupRepository("schema-repo")
		repo.Spec.S3.Endpoint = "s3.example.com"
		err := k8sClient.Create(ctx, repo)
		Expect(apierrors.IsInvalid(err)).To(BeTrue(), "expected Invalid, got %v", err)
	})
})
