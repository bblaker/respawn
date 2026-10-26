package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// BackupDeletionPolicy controls what happens to a snapshot when its GameServerBackup is deleted.
// +kubebuilder:validation:Enum=Retain;Delete
type BackupDeletionPolicy string

// Supported backup deletion policies.
const (
	// BackupRetain keeps the snapshot in the repository.
	BackupRetain BackupDeletionPolicy = "Retain"
	// BackupDelete forgets the snapshot in the repository.
	BackupDelete BackupDeletionPolicy = "Delete"
)

// GameServerBackupSpec defines the desired state of a GameServerBackup.
// It is immutable after creation.
type GameServerBackupSpec struct {
	// gameServerName is the GameServer in the same namespace to back up.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +required
	GameServerName string `json:"gameServerName"`

	// deletionPolicy controls whether deleting this object also forgets the snapshot.
	// +kubebuilder:default=Retain
	// +optional
	DeletionPolicy BackupDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// BackupPhase is a summary of where a backup is in its lifecycle.
// +kubebuilder:validation:Enum=Pending;Quiescing;Running;Succeeded;Failed
type BackupPhase string

// Backup phases.
const (
	BackupPhasePending   BackupPhase = "Pending"
	BackupPhaseQuiescing BackupPhase = "Quiescing"
	BackupPhaseRunning   BackupPhase = "Running"
	BackupPhaseSucceeded BackupPhase = "Succeeded"
	BackupPhaseFailed    BackupPhase = "Failed"
)

// GameServerBackupStatus defines the observed state of a GameServerBackup.
type GameServerBackupStatus struct {
	// observedGeneration is the spec generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// phase summarizes the lifecycle state.
	// +optional
	Phase BackupPhase `json:"phase,omitempty"`

	// quiesced is true while world saving is paused on the game server for this backup.
	// The operator uses it to resume saving after a crash.
	// +optional
	Quiesced bool `json:"quiesced,omitempty"`

	// snapshotID is the restic snapshot ID once the backup has succeeded.
	// +optional
	SnapshotID string `json:"snapshotID,omitempty"`

	// startTime is when the backup started.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// completionTime is when the backup finished, successfully or not.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// sizeBytes is the total size of the backed-up data.
	// +optional
	SizeBytes int64 `json:"sizeBytes,omitempty"`

	// addedBytes is the new data this snapshot added to the repository after deduplication.
	// +optional
	AddedBytes int64 `json:"addedBytes,omitempty"`

	// conditions are the detailed state of the backup.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=respawn
// +kubebuilder:printcolumn:name="Server",type=string,JSONPath=`.spec.gameServerName`
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Snapshot",type=string,JSONPath=`.status.snapshotID`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// GameServerBackup is a single backup of a GameServer's world data.
type GameServerBackup struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of GameServerBackup
	// +required
	Spec GameServerBackupSpec `json:"spec"`

	// status defines the observed state of GameServerBackup
	// +optional
	Status GameServerBackupStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// GameServerBackupList contains a list of GameServerBackup
type GameServerBackupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []GameServerBackup `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &GameServerBackup{}, &GameServerBackupList{})
		return nil
	})
}
