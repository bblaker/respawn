package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GameServerSpec defines the desired state of a GameServer.
//
// Exactly one game block (minecraft or factorio) must be set.
type GameServerSpec struct {
	// minecraft configures a Minecraft Java Edition server.
	// +optional
	Minecraft *MinecraftSpec `json:"minecraft,omitempty"`

	// factorio configures a Factorio headless server.
	// +optional
	Factorio *FactorioSpec `json:"factorio,omitempty"`

	// image overrides the game server image respawn pins for the selected game.
	// Leave empty to use the default.
	// +optional
	Image string `json:"image,omitempty"`

	// storage configures the volume that holds world data.
	// +required
	Storage StorageSpec `json:"storage"`

	// power configures when the server runs.
	// +kubebuilder:default={}
	// +optional
	Power PowerSpec `json:"power,omitzero"`

	// backup configures scheduled and on-delete backups.
	// Leave unset to disable backups.
	// +optional
	Backup *BackupSpec `json:"backup,omitempty"`

	// service configures how players reach the server.
	// +kubebuilder:default={}
	// +optional
	Service ServiceSpec `json:"service,omitzero"`

	// pod holds a constrained set of scheduling and resource overrides for the game pod.
	// +optional
	Pod PodOverrides `json:"pod,omitzero"`
}

// MinecraftType is the server software distribution, passed to the itzg/minecraft-server image.
// +kubebuilder:validation:Enum=VANILLA;PAPER;PURPUR;FABRIC;FORGE;NEOFORGE
type MinecraftType string

// Supported Minecraft server types.
const (
	MinecraftTypeVanilla  MinecraftType = "VANILLA"
	MinecraftTypePaper    MinecraftType = "PAPER"
	MinecraftTypePurpur   MinecraftType = "PURPUR"
	MinecraftTypeFabric   MinecraftType = "FABRIC"
	MinecraftTypeForge    MinecraftType = "FORGE"
	MinecraftTypeNeoForge MinecraftType = "NEOFORGE"
)

// MinecraftSpec configures a Minecraft Java Edition server.
type MinecraftSpec struct {
	// eula records that you accept the Minecraft EULA (https://www.minecraft.net/eula).
	// respawn never accepts it on your behalf; the server will not start unless this is true.
	// +required
	EULA bool `json:"eula"`

	// version is the Minecraft version, for example "1.21.4".
	// It is required rather than defaulted to the latest release, because
	// upgrading a world is irreversible.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +required
	Version string `json:"version"`

	// type is the server software distribution.
	// +kubebuilder:default=VANILLA
	// +optional
	Type MinecraftType `json:"type,omitempty"`

	// memory is the JVM heap size. The container memory limit adds overhead on top of it.
	// +kubebuilder:default="2Gi"
	// +optional
	Memory *resource.Quantity `json:"memory,omitempty"`

	// maxPlayers is the maximum number of players that can be online at once.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1000
	// +kubebuilder:default=20
	// +optional
	MaxPlayers int32 `json:"maxPlayers,omitempty"`

	// motd is the message shown in the multiplayer server list while the server is running.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	MOTD string `json:"motd,omitempty"`

	// env passes extra environment variables to the itzg/minecraft-server image.
	// Keys that respawn manages itself (such as EULA and RCON settings) are rejected.
	// +optional
	Env map[string]string `json:"env,omitempty"`
}

// FactorioSpec configures a Factorio headless server.
type FactorioSpec struct {
	// version is the Factorio version or image tag, for example "2.0" or "stable".
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +required
	Version string `json:"version"`

	// saveName is the name of the save file to load, without the .zip extension.
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9._-]+$`
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:default=world
	// +optional
	SaveName string `json:"saveName,omitempty"`

	// serverSettingsRef names a ConfigMap in the same namespace with a
	// server-settings.json key.
	// +optional
	ServerSettingsRef *corev1.LocalObjectReference `json:"serverSettingsRef,omitempty"`
}

// ReclaimPolicy controls what happens to world data when its GameServer is deleted.
// +kubebuilder:validation:Enum=Retain;Delete
type ReclaimPolicy string

// Supported reclaim policies.
const (
	// ReclaimRetain keeps the PersistentVolumeClaim and labels it as orphaned.
	ReclaimRetain ReclaimPolicy = "Retain"
	// ReclaimDelete deletes the PersistentVolumeClaim.
	ReclaimDelete ReclaimPolicy = "Delete"
)

// StorageSpec configures the volume that holds world data.
type StorageSpec struct {
	// size is the requested size of the world data volume.
	// +required
	Size resource.Quantity `json:"size"`

	// storageClassName is the StorageClass for the volume. Leave unset to use
	// the cluster default.
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty"`

	// reclaimPolicy controls what happens to the volume when the GameServer is deleted.
	// +kubebuilder:default=Retain
	// +optional
	ReclaimPolicy ReclaimPolicy `json:"reclaimPolicy,omitempty"`

	// restoreFrom initializes a new volume from a backup. It is only honored when
	// the volume is first created.
	// +optional
	RestoreFrom *RestoreSource `json:"restoreFrom,omitempty"`
}

// RestoreSource identifies a backup to restore from. Set either backupName,
// or snapshotID together with repositoryRef.
type RestoreSource struct {
	// backupName is a GameServerBackup in the same namespace.
	// +kubebuilder:validation:MaxLength=253
	// +optional
	BackupName string `json:"backupName,omitempty"`

	// snapshotID is a restic snapshot ID in repositoryRef.
	// +kubebuilder:validation:Pattern=`^[0-9a-f]{8,64}$`
	// +optional
	SnapshotID string `json:"snapshotID,omitempty"`

	// repositoryRef names the BackupRepository holding snapshotID.
	// +optional
	RepositoryRef *corev1.LocalObjectReference `json:"repositoryRef,omitempty"`
}

// PowerMode controls when a GameServer runs.
// +kubebuilder:validation:Enum=AlwaysOn;Auto;Off
type PowerMode string

// Supported power modes.
const (
	// PowerAlwaysOn keeps the server running.
	PowerAlwaysOn PowerMode = "AlwaysOn"
	// PowerAuto stops the server when idle and wakes it when a player connects.
	PowerAuto PowerMode = "Auto"
	// PowerOff stops the server and does not wake it.
	PowerOff PowerMode = "Off"
)

// PowerSpec configures when the server runs.
type PowerSpec struct {
	// mode controls when the server runs.
	// +kubebuilder:default=AlwaysOn
	// +optional
	Mode PowerMode `json:"mode,omitempty"`

	// idleTimeout is how long the server must have no players before it sleeps.
	// Only used when mode is Auto.
	// +kubebuilder:default="15m"
	// +optional
	IdleTimeout *metav1.Duration `json:"idleTimeout,omitempty"`

	// pollInterval is how often the operator checks the player count.
	// +kubebuilder:default="30s"
	// +optional
	PollInterval *metav1.Duration `json:"pollInterval,omitempty"`

	// startTimeout is how long the server may take to become ready before it is marked Failed.
	// +kubebuilder:default="10m"
	// +optional
	StartTimeout *metav1.Duration `json:"startTimeout,omitempty"`

	// wakeCooldown ignores repeated wake requests within this window.
	// +kubebuilder:default="60s"
	// +optional
	WakeCooldown *metav1.Duration `json:"wakeCooldown,omitempty"`

	// sleepingMotd is the server list message shown while the server sleeps.
	// +kubebuilder:validation:MaxLength=256
	// +optional
	SleepingMOTD string `json:"sleepingMotd,omitempty"`
}

// BackupSpec configures scheduled and on-delete backups.
type BackupSpec struct {
	// repositoryRef names the BackupRepository in the same namespace that backups are written to.
	// +required
	RepositoryRef corev1.LocalObjectReference `json:"repositoryRef"`

	// schedule is a standard 5-field cron expression. Leave empty for manual backups only.
	// +kubebuilder:validation:MaxLength=128
	// +optional
	Schedule string `json:"schedule,omitempty"`

	// retention controls which scheduled snapshots are kept.
	// +optional
	Retention *RetentionPolicy `json:"retention,omitempty"`

	// beforeSleep takes a backup each time the server goes to sleep, at most once per minInterval.
	// +optional
	BeforeSleep bool `json:"beforeSleep,omitempty"`

	// minInterval throttles beforeSleep backups.
	// +kubebuilder:default="1h"
	// +optional
	MinInterval *metav1.Duration `json:"minInterval,omitempty"`

	// onDelete takes a final backup before the GameServer is deleted.
	// +kubebuilder:default=true
	// +optional
	OnDelete *bool `json:"onDelete,omitempty"`
}

// RetentionPolicy maps to restic forget --keep-* flags. Unset fields keep nothing for that bucket.
type RetentionPolicy struct {
	// keepLast keeps the most recent N snapshots.
	// +kubebuilder:validation:Minimum=0
	// +optional
	KeepLast *int32 `json:"keepLast,omitempty"`

	// keepDaily keeps the most recent snapshot for each of the last N days.
	// +kubebuilder:validation:Minimum=0
	// +optional
	KeepDaily *int32 `json:"keepDaily,omitempty"`

	// keepWeekly keeps the most recent snapshot for each of the last N weeks.
	// +kubebuilder:validation:Minimum=0
	// +optional
	KeepWeekly *int32 `json:"keepWeekly,omitempty"`

	// keepMonthly keeps the most recent snapshot for each of the last N months.
	// +kubebuilder:validation:Minimum=0
	// +optional
	KeepMonthly *int32 `json:"keepMonthly,omitempty"`
}

// ServiceSpec configures the Service that players connect to. Each GameServer
// gets its own Service.
type ServiceSpec struct {
	// type is the Kubernetes Service type.
	// +kubebuilder:validation:Enum=LoadBalancer;NodePort;ClusterIP
	// +kubebuilder:default=LoadBalancer
	// +optional
	Type corev1.ServiceType `json:"type,omitempty"`

	// port is the external port. It defaults to the game's standard port.
	// With k3s ServiceLB, every server on a node needs a different port.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port *int32 `json:"port,omitempty"`

	// annotations are added to the Service, for example
	// metallb.universe.tf/loadBalancerIPs to request a fixed IP.
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// PodOverrides is a deliberately narrow set of pod settings. A full PodTemplate
// is not exposed, so a GameServer cannot request privileged or host access.
type PodOverrides struct {
	// nodeSelector constrains the game pod to matching nodes.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// tolerations for the game pod.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// affinity for the game pod.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// resources for the game container. Memory is derived from the game block
	// unless set here.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`
}

// GameServerPhase is a summary of where a GameServer is in its lifecycle.
// Conditions carry the detail.
// +kubebuilder:validation:Enum=Pending;Starting;Running;Stopping;Sleeping;Stopped;Failed
type GameServerPhase string

// GameServer phases.
const (
	PhasePending  GameServerPhase = "Pending"
	PhaseStarting GameServerPhase = "Starting"
	PhaseRunning  GameServerPhase = "Running"
	PhaseStopping GameServerPhase = "Stopping"
	PhaseSleeping GameServerPhase = "Sleeping"
	PhaseStopped  GameServerPhase = "Stopped"
	PhaseFailed   GameServerPhase = "Failed"
)

// BackupReference points at a completed backup.
type BackupReference struct {
	// name of the GameServerBackup.
	Name string `json:"name"`

	// time the backup completed.
	// +optional
	Time *metav1.Time `json:"time,omitempty"`

	// snapshotID is the restic snapshot ID.
	// +optional
	SnapshotID string `json:"snapshotID,omitempty"`
}

// GameServerStatus defines the observed state of a GameServer.
type GameServerStatus struct {
	// observedGeneration is the spec generation this status was computed from.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// game is the game this server runs, for example "minecraft".
	// +optional
	Game string `json:"game,omitempty"`

	// phase summarizes the lifecycle state.
	// +optional
	Phase GameServerPhase `json:"phase,omitempty"`

	// playersOnline is the last observed player count. Unset when unknown.
	// +optional
	PlayersOnline *int32 `json:"playersOnline,omitempty"`

	// maxPlayers is the server's player limit.
	// +optional
	MaxPlayers int32 `json:"maxPlayers,omitempty"`

	// endpoint is the address players connect to, as host:port.
	// +optional
	Endpoint string `json:"endpoint,omitempty"`

	// idleSince is when the player count was first observed to be zero. Cleared when players join.
	// +optional
	IdleSince *metav1.Time `json:"idleSince,omitempty"`

	// readySince is when the server last became ready.
	// +optional
	ReadySince *metav1.Time `json:"readySince,omitempty"`

	// lastWakeTime is when the last accepted wake request was handled.
	// +optional
	LastWakeTime *metav1.Time `json:"lastWakeTime,omitempty"`

	// lastBackup is the most recent successful backup.
	// +optional
	LastBackup *BackupReference `json:"lastBackup,omitempty"`

	// nextBackupTime is when the next scheduled backup is due.
	// +optional
	NextBackupTime *metav1.Time `json:"nextBackupTime,omitempty"`

	// conditions are the detailed state of the GameServer: Ready, Available,
	// Progressing, Degraded, and BackupHealthy.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories=respawn
// +kubebuilder:printcolumn:name="Game",type=string,JSONPath=`.status.game`
// +kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Players",type=integer,JSONPath=`.status.playersOnline`
// +kubebuilder:printcolumn:name="Endpoint",type=string,JSONPath=`.status.endpoint`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// GameServer is a long-lived game server with persistent world data,
// backups, and optional scale-to-zero.
type GameServer struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of GameServer
	// +required
	Spec GameServerSpec `json:"spec"`

	// status defines the observed state of GameServer
	// +optional
	Status GameServerStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// GameServerList contains a list of GameServer
type GameServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []GameServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &GameServer{}, &GameServerList{})
		return nil
	})
}
