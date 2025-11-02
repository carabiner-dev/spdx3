package software

import "github.com/carabiner-dev/spdx3/types"

// ContentIdentifierType vocabulary defines valid content identifier types
type ContentIdentifierType string

var ContentIdentifierTypes = types.Vocabulary[ContentIdentifierType]{
	ContentIdentifierTypeGitoid,
	ContentIdentifierTypeSwhid,
}

const (
	// ContentIdentifierTypeGitoid represents a Git Object ID - a unique hash of a binary artifact
	ContentIdentifierTypeGitoid ContentIdentifierType = "gitoid"
	// ContentIdentifierTypeSwhid represents a SoftWare Hash IDentifier - a persistent intrinsic identifier for digital artifacts
	ContentIdentifierTypeSwhid ContentIdentifierType = "swhid"
)

// FileKindType vocabulary defines valid file kind types
type FileKindType string

var FileKindTypes = types.Vocabulary[FileKindType]{
	FileKindTypeDirectory,
	FileKindTypeFile,
}

const (
	// FileKindTypeDirectory represents a directory and all content stored in that directory
	FileKindTypeDirectory FileKindType = "directory"
	// FileKindTypeFile represents a single file (default)
	FileKindTypeFile FileKindType = "file"
)

// SbomType vocabulary defines valid SBOM document types
type SbomType string

var SbomTypes = types.Vocabulary[SbomType]{
	SbomTypeAnalyzed,
	SbomTypeBuild,
	SbomTypeDeployed,
	SbomTypeDesign,
	SbomTypeRuntime,
	SbomTypeSource,
}

const (
	// SbomTypeAnalyzed represents an SBOM generated through analysis of artifacts after its build
	SbomTypeAnalyzed SbomType = "analyzed"
	// SbomTypeBuild represents an SBOM generated as part of the build process
	SbomTypeBuild SbomType = "build"
	// SbomTypeDeployed represents an SBOM providing inventory of software present on a system
	SbomTypeDeployed SbomType = "deployed"
	// SbomTypeDesign represents an SBOM of intended, planned software project or product
	SbomTypeDesign SbomType = "design"
	// SbomTypeRuntime represents an SBOM generated through instrumenting the running system
	SbomTypeRuntime SbomType = "runtime"
	// SbomTypeSource represents an SBOM created directly from the development environment
	SbomTypeSource SbomType = "source"
)

// SoftwarePurpose vocabulary defines valid software purpose types
type SoftwarePurpose string

var SoftwarePurposes = types.Vocabulary[SoftwarePurpose]{
	SoftwarePurposeApplication,
	SoftwarePurposeArchive,
	SoftwarePurposeBom,
	SoftwarePurposeConfiguration,
	SoftwarePurposeContainer,
	SoftwarePurposeData,
	SoftwarePurposeDevice,
	SoftwarePurposeDeviceDriver,
	SoftwarePurposeDiskImage,
	SoftwarePurposeDocumentation,
	SoftwarePurposeEvidence,
	SoftwarePurposeExecutable,
	SoftwarePurposeFile,
	SoftwarePurposeFilesystemImage,
	SoftwarePurposeFirmware,
	SoftwarePurposeFramework,
	SoftwarePurposeInstall,
	SoftwarePurposeLibrary,
	SoftwarePurposeManifest,
	SoftwarePurposeModel,
	SoftwarePurposeModule,
	SoftwarePurposeOperatingSystem,
	SoftwarePurposeOther,
	SoftwarePurposePatch,
	SoftwarePurposePlatform,
	SoftwarePurposeRequirement,
	SoftwarePurposeSource,
	SoftwarePurposeSpecification,
	SoftwarePurposeTest,
}

const (
	// SoftwarePurposeApplication indicates the element is an application
	SoftwarePurposeApplication SoftwarePurpose = "application"
	// SoftwarePurposeArchive indicates the element is an archive
	SoftwarePurposeArchive SoftwarePurpose = "archive"
	// SoftwarePurposeBom indicates the element is a bill of materials
	SoftwarePurposeBom SoftwarePurpose = "bom"
	// SoftwarePurposeConfiguration indicates the element is configuration
	SoftwarePurposeConfiguration SoftwarePurpose = "configuration"
	// SoftwarePurposeContainer indicates the element is a container
	SoftwarePurposeContainer SoftwarePurpose = "container"
	// SoftwarePurposeData indicates the element is data
	SoftwarePurposeData SoftwarePurpose = "data"
	// SoftwarePurposeDevice indicates the element is a device
	SoftwarePurposeDevice SoftwarePurpose = "device"
	// SoftwarePurposeDeviceDriver indicates the element is a device driver
	SoftwarePurposeDeviceDriver SoftwarePurpose = "deviceDriver"
	// SoftwarePurposeDiskImage indicates the element is a disk image
	SoftwarePurposeDiskImage SoftwarePurpose = "diskImage"
	// SoftwarePurposeDocumentation indicates the element is documentation
	SoftwarePurposeDocumentation SoftwarePurpose = "documentation"
	// SoftwarePurposeEvidence indicates the element is evidence
	SoftwarePurposeEvidence SoftwarePurpose = "evidence"
	// SoftwarePurposeExecutable indicates the element is an executable
	SoftwarePurposeExecutable SoftwarePurpose = "executable"
	// SoftwarePurposeFile indicates the element is a file
	SoftwarePurposeFile SoftwarePurpose = "file"
	// SoftwarePurposeFilesystemImage indicates the element is a filesystem image
	SoftwarePurposeFilesystemImage SoftwarePurpose = "filesystemImage"
	// SoftwarePurposeFirmware indicates the element is firmware
	SoftwarePurposeFirmware SoftwarePurpose = "firmware"
	// SoftwarePurposeFramework indicates the element is a framework
	SoftwarePurposeFramework SoftwarePurpose = "framework"
	// SoftwarePurposeInstall indicates the element is an installer
	SoftwarePurposeInstall SoftwarePurpose = "install"
	// SoftwarePurposeLibrary indicates the element is a library
	SoftwarePurposeLibrary SoftwarePurpose = "library"
	// SoftwarePurposeManifest indicates the element is a manifest
	SoftwarePurposeManifest SoftwarePurpose = "manifest"
	// SoftwarePurposeModel indicates the element is a model
	SoftwarePurposeModel SoftwarePurpose = "model"
	// SoftwarePurposeModule indicates the element is a module
	SoftwarePurposeModule SoftwarePurpose = "module"
	// SoftwarePurposeOperatingSystem indicates the element is an operating system
	SoftwarePurposeOperatingSystem SoftwarePurpose = "operatingSystem"
	// SoftwarePurposeOther indicates the element has another purpose
	SoftwarePurposeOther SoftwarePurpose = "other"
	// SoftwarePurposePatch indicates the element is a patch
	SoftwarePurposePatch SoftwarePurpose = "patch"
	// SoftwarePurposePlatform indicates the element is a platform
	SoftwarePurposePlatform SoftwarePurpose = "platform"
	// SoftwarePurposeRequirement indicates the element is a requirement
	SoftwarePurposeRequirement SoftwarePurpose = "requirement"
	// SoftwarePurposeSource indicates the element is source code
	SoftwarePurposeSource SoftwarePurpose = "source"
	// SoftwarePurposeSpecification indicates the element is a specification
	SoftwarePurposeSpecification SoftwarePurpose = "specification"
	// SoftwarePurposeTest indicates the element is a test
	SoftwarePurposeTest SoftwarePurpose = "test"
)
