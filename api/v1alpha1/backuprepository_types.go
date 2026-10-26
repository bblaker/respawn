package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// BackupRepositorySpec defines the desired state of a BackupRepository.
type BackupRepositorySpec struct {
	// s3 configures an S3-compatible bucket.
	// +required
	S3 S3Spec `json:"s3"`

	// credentialsSecretRef names a Secret in the same namespace with
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY keys.
	// +required
	CredentialsSecretRef corev1.LocalObjectReference `json:"credentialsSecretRef"`

	// encryptionSecretRef names a Secret in the same namespace with a
	// RESTIC_PASSWORD key. Backups cannot be restored without it.
	// +required
	EncryptionSecretRef corev1.LocalObjectReference `json:"encryptionSecretRef"`
}

// S3Spec configures an S3-compatible bucket.
type S3Spec struct {
	// endpoint is the S3 API URL, for example https://s3.us-west-004.backblazeb2.com.
	// +kubebuilder:validation:Pattern=`^https?://`
	// +kubebuilder:validation:MaxLength=2048
	// +required
	Endpoint string `json:"endpoint"`

	// bucket is the bucket name.
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=63
	// +required
	Bucket string `json:"bucket"`

	// prefix is prepended to every repository path in the bucket.
	// Each GameServer gets its own path under it.
	// +kubebuilder:validation:MaxLength=512
	// +optional
	Prefix string `json:"prefix,omitempty"`

	// region is the bucket region, if the provider needs one.
	// +kubebuilder:validation:MaxLength=64
	// +optional
	Region string `json:"region,omitempty"`
}

// BackupRepositoryStatus defines the observed state of a BackupRepository.
type BackupRepositoryStatus struct {
	// observedGeneration is the spec generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// conditions are the detailed state of the repository. Ready is true when
	// the last connectivity check succeeded.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=respawn
// +kubebuilder:printcolumn:name="Bucket",type=string,JSONPath=`.spec.s3.bucket`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// BackupRepository is an S3-compatible location where restic backups are stored.
type BackupRepository struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of BackupRepository
	// +required
	Spec BackupRepositorySpec `json:"spec"`

	// status defines the observed state of BackupRepository
	// +optional
	Status BackupRepositoryStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// BackupRepositoryList contains a list of BackupRepository
type BackupRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []BackupRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &BackupRepository{}, &BackupRepositoryList{})
		return nil
	})
}
