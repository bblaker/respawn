package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	respawnv1alpha1 "github.com/bblaker/respawn/api/v1alpha1"
)

// testNamespace is where envtest specs create their objects.
const testNamespace = "default"

// newGameServer returns the smallest valid Minecraft GameServer.
func newGameServer(name string) *respawnv1alpha1.GameServer {
	return &respawnv1alpha1.GameServer{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace},
		Spec: respawnv1alpha1.GameServerSpec{
			Minecraft: &respawnv1alpha1.MinecraftSpec{EULA: true, Version: "1.21.4"},
			Storage:   respawnv1alpha1.StorageSpec{Size: resource.MustParse("1Gi")},
		},
	}
}

// newGameServerBackup returns a valid GameServerBackup for the named server.
func newGameServerBackup(name, gameServer string) *respawnv1alpha1.GameServerBackup {
	return &respawnv1alpha1.GameServerBackup{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace},
		Spec:       respawnv1alpha1.GameServerBackupSpec{GameServerName: gameServer},
	}
}

// newBackupRepository returns a valid BackupRepository.
func newBackupRepository(name string) *respawnv1alpha1.BackupRepository {
	return &respawnv1alpha1.BackupRepository{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: testNamespace},
		Spec: respawnv1alpha1.BackupRepositorySpec{
			S3: respawnv1alpha1.S3Spec{
				Endpoint: "http://minio.minio.svc:9000",
				Bucket:   "respawn-test",
			},
			CredentialsSecretRef: corev1.LocalObjectReference{Name: "s3-credentials"},
			EncryptionSecretRef:  corev1.LocalObjectReference{Name: "restic-password"},
		},
	}
}
