// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package core

import "github.com/carabiner-dev/spdx3/types"

// AnnotationType vocabulary defines valid annotation types
type AnnotationType string

var AnnotationTypes = types.Vocabulary[AnnotationType]{
	AnnotationTypeOther,
	AnnotationTypeReview,
}

const (
	// AnnotationTypeOther for supplementary information outside formal review
	AnnotationTypeOther AnnotationType = "other"
	// AnnotationTypeReview for Element evaluation by a reviewer
	AnnotationTypeReview AnnotationType = "review"
)

// ExternalIdentifierType vocabulary defines valid external identifier types
type ExternalIdentifierType string

var ExternalIdentifierTypes = types.Vocabulary[ExternalIdentifierType]{
	ExternalIdentifierTypeCpe22,
	ExternalIdentifierTypeCpe23,
	ExternalIdentifierTypeCve,
	ExternalIdentifierTypeEmail,
	ExternalIdentifierTypeGitoid,
	ExternalIdentifierTypeOther,
	ExternalIdentifierTypePackageUrl,
	ExternalIdentifierTypeSecurityOther,
	ExternalIdentifierTypeSwhid,
	ExternalIdentifierTypeSwid,
	ExternalIdentifierTypeUrlScheme,
}

const (
	ExternalIdentifierTypeCpe22         ExternalIdentifierType = "cpe22"
	ExternalIdentifierTypeCpe23         ExternalIdentifierType = "cpe23"
	ExternalIdentifierTypeCve           ExternalIdentifierType = "cve"
	ExternalIdentifierTypeEmail         ExternalIdentifierType = "email"
	ExternalIdentifierTypeGitoid        ExternalIdentifierType = "gitoid"
	ExternalIdentifierTypeOther         ExternalIdentifierType = "other"
	ExternalIdentifierTypePackageUrl    ExternalIdentifierType = "packageUrl"
	ExternalIdentifierTypeSecurityOther ExternalIdentifierType = "securityOther"
	ExternalIdentifierTypeSwhid         ExternalIdentifierType = "swhid"
	ExternalIdentifierTypeSwid          ExternalIdentifierType = "swid"
	ExternalIdentifierTypeUrlScheme     ExternalIdentifierType = "urlScheme"
)

// ExternalRefType vocabulary defines valid external reference types
type ExternalRefType string

var ExternalRefTypes = types.Vocabulary[ExternalRefType]{
	ExternalRefTypeAltDownloadLocation,
	ExternalRefTypeAltWebPage,
	ExternalRefTypeBinaryArtifact,
	ExternalRefTypeBower,
	ExternalRefTypeBuildMeta,
	ExternalRefTypeBuildSystem,
	ExternalRefTypeCertificationReport,
	ExternalRefTypeChat,
	ExternalRefTypeComponentAnalysisReport,
	ExternalRefTypeCwe,
	ExternalRefTypeDocumentation,
	ExternalRefTypeDynamicAnalysisReport,
	ExternalRefTypeEolNotice,
	ExternalRefTypeExportControlAssessment,
	ExternalRefTypeFunding,
	ExternalRefTypeIssueTracker,
	ExternalRefTypeLicense,
	ExternalRefTypeMailingList,
	ExternalRefTypeMavenCentral,
	ExternalRefTypeMetrics,
	ExternalRefTypeNpm,
	ExternalRefTypeNuget,
	ExternalRefTypeOther,
	ExternalRefTypePrivacyAssessment,
	ExternalRefTypeProductMetadata,
	ExternalRefTypePurchaseOrder,
	ExternalRefTypeQualityAssessmentReport,
	ExternalRefTypeReleaseHistory,
	ExternalRefTypeReleaseNotes,
	ExternalRefTypeRiskAssessment,
	ExternalRefTypeRuntimeAnalysisReport,
	ExternalRefTypeSecurityAdversaryModel,
	ExternalRefTypeSecurityAdvisory,
	ExternalRefTypeSecurityFix,
	ExternalRefTypeSecurityOther,
	ExternalRefTypeSecurityPenTestReport,
	ExternalRefTypeSecurityPolicy,
	ExternalRefTypeSecurityThreatModel,
	ExternalRefTypeSecureSoftwareAttestation,
	ExternalRefTypeSocialMedia,
	ExternalRefTypeSourceArtifact,
	ExternalRefTypeStaticAnalysisReport,
	ExternalRefTypeSupport,
	ExternalRefTypeVcs,
	ExternalRefTypeVulnerabilityDisclosureReport,
	ExternalRefTypeVulnerabilityExploitabilityAssessment,
}

const (
	ExternalRefTypeAltDownloadLocation                   ExternalRefType = "altDownloadLocation"
	ExternalRefTypeAltWebPage                            ExternalRefType = "altWebPage"
	ExternalRefTypeBinaryArtifact                        ExternalRefType = "binaryArtifact"
	ExternalRefTypeBower                                 ExternalRefType = "bower"
	ExternalRefTypeBuildMeta                             ExternalRefType = "buildMeta"
	ExternalRefTypeBuildSystem                           ExternalRefType = "buildSystem"
	ExternalRefTypeCertificationReport                   ExternalRefType = "certificationReport"
	ExternalRefTypeChat                                  ExternalRefType = "chat"
	ExternalRefTypeComponentAnalysisReport               ExternalRefType = "componentAnalysisReport"
	ExternalRefTypeCwe                                   ExternalRefType = "cwe"
	ExternalRefTypeDocumentation                         ExternalRefType = "documentation"
	ExternalRefTypeDynamicAnalysisReport                 ExternalRefType = "dynamicAnalysisReport"
	ExternalRefTypeEolNotice                             ExternalRefType = "eolNotice"
	ExternalRefTypeExportControlAssessment               ExternalRefType = "exportControlAssessment"
	ExternalRefTypeFunding                               ExternalRefType = "funding"
	ExternalRefTypeIssueTracker                          ExternalRefType = "issueTracker"
	ExternalRefTypeLicense                               ExternalRefType = "license"
	ExternalRefTypeMailingList                           ExternalRefType = "mailingList"
	ExternalRefTypeMavenCentral                          ExternalRefType = "mavenCentral"
	ExternalRefTypeMetrics                               ExternalRefType = "metrics"
	ExternalRefTypeNpm                                   ExternalRefType = "npm"
	ExternalRefTypeNuget                                 ExternalRefType = "nuget"
	ExternalRefTypeOther                                 ExternalRefType = "other"
	ExternalRefTypePrivacyAssessment                     ExternalRefType = "privacyAssessment"
	ExternalRefTypeProductMetadata                       ExternalRefType = "productMetadata"
	ExternalRefTypePurchaseOrder                         ExternalRefType = "purchaseOrder"
	ExternalRefTypeQualityAssessmentReport               ExternalRefType = "qualityAssessmentReport"
	ExternalRefTypeReleaseHistory                        ExternalRefType = "releaseHistory"
	ExternalRefTypeReleaseNotes                          ExternalRefType = "releaseNotes"
	ExternalRefTypeRiskAssessment                        ExternalRefType = "riskAssessment"
	ExternalRefTypeRuntimeAnalysisReport                 ExternalRefType = "runtimeAnalysisReport"
	ExternalRefTypeSecurityAdversaryModel                ExternalRefType = "securityAdversaryModel"
	ExternalRefTypeSecurityAdvisory                      ExternalRefType = "securityAdvisory"
	ExternalRefTypeSecurityFix                           ExternalRefType = "securityFix"
	ExternalRefTypeSecurityOther                         ExternalRefType = "securityOther"
	ExternalRefTypeSecurityPenTestReport                 ExternalRefType = "securityPenTestReport"
	ExternalRefTypeSecurityPolicy                        ExternalRefType = "securityPolicy"
	ExternalRefTypeSecurityThreatModel                   ExternalRefType = "securityThreatModel"
	ExternalRefTypeSecureSoftwareAttestation             ExternalRefType = "secureSoftwareAttestation"
	ExternalRefTypeSocialMedia                           ExternalRefType = "socialMedia"
	ExternalRefTypeSourceArtifact                        ExternalRefType = "sourceArtifact"
	ExternalRefTypeStaticAnalysisReport                  ExternalRefType = "staticAnalysisReport"
	ExternalRefTypeSupport                               ExternalRefType = "support"
	ExternalRefTypeVcs                                   ExternalRefType = "vcs"
	ExternalRefTypeVulnerabilityDisclosureReport         ExternalRefType = "vulnerabilityDisclosureReport"
	ExternalRefTypeVulnerabilityExploitabilityAssessment ExternalRefType = "vulnerabilityExploitabilityAssessment"
)

// HashAlgorithm vocabulary defines valid hash algorithms
type HashAlgorithm string

var HashAlgorithms = types.Vocabulary[HashAlgorithm]{
	HashAlgorithmAdler32,
	HashAlgorithmBlake2b256,
	HashAlgorithmBlake2b384,
	HashAlgorithmBlake2b512,
	HashAlgorithmBlake3,
	HashAlgorithmCrystalsDilithium,
	HashAlgorithmCrystalsKyber,
	HashAlgorithmFalcon,
	HashAlgorithmMd2,
	HashAlgorithmMd4,
	HashAlgorithmMd5,
	HashAlgorithmMd6,
	HashAlgorithmOther,
	HashAlgorithmSha1,
	HashAlgorithmSha224,
	HashAlgorithmSha256,
	HashAlgorithmSha384,
	HashAlgorithmSha3_224,
	HashAlgorithmSha3_256,
	HashAlgorithmSha3_384,
	HashAlgorithmSha3_512,
	HashAlgorithmSha512,
}

const (
	HashAlgorithmAdler32           HashAlgorithm = "adler32"
	HashAlgorithmBlake2b256        HashAlgorithm = "blake2b256"
	HashAlgorithmBlake2b384        HashAlgorithm = "blake2b384"
	HashAlgorithmBlake2b512        HashAlgorithm = "blake2b512"
	HashAlgorithmBlake3            HashAlgorithm = "blake3"
	HashAlgorithmCrystalsDilithium HashAlgorithm = "crystalsDilithium"
	HashAlgorithmCrystalsKyber     HashAlgorithm = "crystalsKyber"
	HashAlgorithmFalcon            HashAlgorithm = "falcon"
	HashAlgorithmMd2               HashAlgorithm = "md2"
	HashAlgorithmMd4               HashAlgorithm = "md4"
	HashAlgorithmMd5               HashAlgorithm = "md5"
	HashAlgorithmMd6               HashAlgorithm = "md6"
	HashAlgorithmOther             HashAlgorithm = "other"
	HashAlgorithmSha1              HashAlgorithm = "sha1"
	HashAlgorithmSha224            HashAlgorithm = "sha224"
	HashAlgorithmSha256            HashAlgorithm = "sha256"
	HashAlgorithmSha384            HashAlgorithm = "sha384"
	HashAlgorithmSha3_224          HashAlgorithm = "sha3_224"
	HashAlgorithmSha3_256          HashAlgorithm = "sha3_256"
	HashAlgorithmSha3_384          HashAlgorithm = "sha3_384"
	HashAlgorithmSha3_512          HashAlgorithm = "sha3_512"
	HashAlgorithmSha512            HashAlgorithm = "sha512"
)

// LifecycleScopeType vocabulary defines valid lifecycle scope types
type LifecycleScopeType string

var LifecycleScopeTypes = types.Vocabulary[LifecycleScopeType]{
	LifecycleScopeTypeBuild,
	LifecycleScopeTypeDesign,
	LifecycleScopeTypeDevelopment,
	LifecycleScopeTypeOther,
	LifecycleScopeTypeRuntime,
	LifecycleScopeTypeTest,
}

const (
	LifecycleScopeTypeBuild       LifecycleScopeType = "build"
	LifecycleScopeTypeDesign      LifecycleScopeType = "design"
	LifecycleScopeTypeDevelopment LifecycleScopeType = "development"
	LifecycleScopeTypeOther       LifecycleScopeType = "other"
	LifecycleScopeTypeRuntime     LifecycleScopeType = "runtime"
	LifecycleScopeTypeTest        LifecycleScopeType = "test"
)

// PresenceType vocabulary defines valid presence types
type PresenceType string

var PresenceTypes = types.Vocabulary[PresenceType]{
	PresenceTypeNo,
	PresenceTypeNoAssertion,
	PresenceTypeYes,
}

const (
	// PresenceTypeNo indicates absence of the field
	PresenceTypeNo PresenceType = "no"
	// PresenceTypeNoAssertion makes no assertion about the field
	PresenceTypeNoAssertion PresenceType = "noAssertion"
	// PresenceTypeYes indicates presence of the field
	PresenceTypeYes PresenceType = "yes"
)

// ProfileIdentifierType vocabulary defines valid profile identifiers
type ProfileIdentifierType string

var ProfileIdentifierTypes = types.Vocabulary[ProfileIdentifierType]{
	ProfileIdentifierTypeAi,
	ProfileIdentifierTypeBuild,
	ProfileIdentifierTypeCore,
	ProfileIdentifierTypeDataset,
	ProfileIdentifierTypeExpandedLicensing,
	ProfileIdentifierTypeExtension,
	ProfileIdentifierTypeLite,
	ProfileIdentifierTypeSecurity,
	ProfileIdentifierTypeSimpleLicensing,
	ProfileIdentifierTypeSoftware,
}

const (
	ProfileIdentifierTypeAi                ProfileIdentifierType = "ai"
	ProfileIdentifierTypeBuild             ProfileIdentifierType = "build"
	ProfileIdentifierTypeCore              ProfileIdentifierType = "core"
	ProfileIdentifierTypeDataset           ProfileIdentifierType = "dataset"
	ProfileIdentifierTypeExpandedLicensing ProfileIdentifierType = "expandedLicensing"
	ProfileIdentifierTypeExtension         ProfileIdentifierType = "extension"
	ProfileIdentifierTypeLite              ProfileIdentifierType = "lite"
	ProfileIdentifierTypeSecurity          ProfileIdentifierType = "security"
	ProfileIdentifierTypeSimpleLicensing   ProfileIdentifierType = "simpleLicensing"
	ProfileIdentifierTypeSoftware          ProfileIdentifierType = "software"
)

// RelationshipCompleteness vocabulary defines valid relationship completeness types
type RelationshipCompleteness string

var RelationshipCompletenesses = types.Vocabulary[RelationshipCompleteness]{
	RelationshipCompletenessComplete,
	RelationshipCompletenessIncomplete,
	RelationshipCompletenessNoAssertion,
}

const (
	// RelationshipCompletenessComplete indicates the relationship is exhaustive
	RelationshipCompletenessComplete RelationshipCompleteness = "complete"
	// RelationshipCompletenessIncomplete indicates the relationship is not exhaustive
	RelationshipCompletenessIncomplete RelationshipCompleteness = "incomplete"
	// RelationshipCompletenessNoAssertion indicates no assertion about completeness
	RelationshipCompletenessNoAssertion RelationshipCompleteness = "noAssertion"
)

// RelationshipType vocabulary defines valid relationship types
type RelationshipType string

var RelationshipTypes = types.Vocabulary[RelationshipType]{
	RelationshipTypeAffects,
	RelationshipTypeAmendedBy,
	RelationshipTypeAncestorOf,
	RelationshipTypeAvailableFrom,
	RelationshipTypeConfigures,
	RelationshipTypeContains,
	RelationshipTypeCoordinatedBy,
	RelationshipTypeCopiedTo,
	RelationshipTypeDelegatedTo,
	RelationshipTypeDependsOn,
	RelationshipTypeDescendantOf,
	RelationshipTypeDescribes,
	RelationshipTypeDoesNotAffect,
	RelationshipTypeExpandsTo,
	RelationshipTypeExploitCreatedBy,
	RelationshipTypeFixedBy,
	RelationshipTypeFixedIn,
	RelationshipTypeFoundBy,
	RelationshipTypeGenerates,
	RelationshipTypeHasAddedFile,
	RelationshipTypeHasAssessmentFor,
	RelationshipTypeHasAssociatedVulnerability,
	RelationshipTypeHasConcludedLicense,
	RelationshipTypeHasDataFile,
	RelationshipTypeHasDeclaredLicense,
	RelationshipTypeHasDeletedFile,
	RelationshipTypeHasDependencyManifest,
	RelationshipTypeHasDistributionArtifact,
	RelationshipTypeHasDocumentation,
	RelationshipTypeHasDynamicLink,
	RelationshipTypeHasEvidence,
	RelationshipTypeHasExample,
	RelationshipTypeHasHost,
	RelationshipTypeHasInput,
	RelationshipTypeHasMetadata,
	RelationshipTypeHasOptionalComponent,
	RelationshipTypeHasOptionalDependency,
	RelationshipTypeHasOutput,
	RelationshipTypeHasPrerequisite,
	RelationshipTypeHasProvidedDependency,
	RelationshipTypeHasRequirement,
	RelationshipTypeHasSpecification,
	RelationshipTypeHasStaticLink,
	RelationshipTypeHasTest,
	RelationshipTypeHasTestCase,
	RelationshipTypeHasVariant,
	RelationshipTypeInvokedBy,
	RelationshipTypeModifiedBy,
	RelationshipTypeOther,
	RelationshipTypePackagedBy,
	RelationshipTypePatchedBy,
	RelationshipTypePublishedBy,
	RelationshipTypeReportedBy,
	RelationshipTypeRepublishedBy,
	RelationshipTypeSerializedInArtifact,
	RelationshipTypeTestedOn,
	RelationshipTypeTrainedOn,
	RelationshipTypeUnderInvestigationFor,
	RelationshipTypeUsesTool,
}

const (
	RelationshipTypeAffects                    RelationshipType = "affects"
	RelationshipTypeAmendedBy                  RelationshipType = "amendedBy"
	RelationshipTypeAncestorOf                 RelationshipType = "ancestorOf"
	RelationshipTypeAvailableFrom              RelationshipType = "availableFrom"
	RelationshipTypeConfigures                 RelationshipType = "configures"
	RelationshipTypeContains                   RelationshipType = "contains"
	RelationshipTypeCoordinatedBy              RelationshipType = "coordinatedBy"
	RelationshipTypeCopiedTo                   RelationshipType = "copiedTo"
	RelationshipTypeDelegatedTo                RelationshipType = "delegatedTo"
	RelationshipTypeDependsOn                  RelationshipType = "dependsOn"
	RelationshipTypeDescendantOf               RelationshipType = "descendantOf"
	RelationshipTypeDescribes                  RelationshipType = "describes"
	RelationshipTypeDoesNotAffect              RelationshipType = "doesNotAffect"
	RelationshipTypeExpandsTo                  RelationshipType = "expandsTo"
	RelationshipTypeExploitCreatedBy           RelationshipType = "exploitCreatedBy"
	RelationshipTypeFixedBy                    RelationshipType = "fixedBy"
	RelationshipTypeFixedIn                    RelationshipType = "fixedIn"
	RelationshipTypeFoundBy                    RelationshipType = "foundBy"
	RelationshipTypeGenerates                  RelationshipType = "generates"
	RelationshipTypeHasAddedFile               RelationshipType = "hasAddedFile"
	RelationshipTypeHasAssessmentFor           RelationshipType = "hasAssessmentFor"
	RelationshipTypeHasAssociatedVulnerability RelationshipType = "hasAssociatedVulnerability"
	RelationshipTypeHasConcludedLicense        RelationshipType = "hasConcludedLicense"
	RelationshipTypeHasDataFile                RelationshipType = "hasDataFile"
	RelationshipTypeHasDeclaredLicense         RelationshipType = "hasDeclaredLicense"
	RelationshipTypeHasDeletedFile             RelationshipType = "hasDeletedFile"
	RelationshipTypeHasDependencyManifest      RelationshipType = "hasDependencyManifest"
	RelationshipTypeHasDistributionArtifact    RelationshipType = "hasDistributionArtifact"
	RelationshipTypeHasDocumentation           RelationshipType = "hasDocumentation"
	RelationshipTypeHasDynamicLink             RelationshipType = "hasDynamicLink"
	RelationshipTypeHasEvidence                RelationshipType = "hasEvidence"
	RelationshipTypeHasExample                 RelationshipType = "hasExample"
	RelationshipTypeHasHost                    RelationshipType = "hasHost"
	RelationshipTypeHasInput                   RelationshipType = "hasInput"
	RelationshipTypeHasMetadata                RelationshipType = "hasMetadata"
	RelationshipTypeHasOptionalComponent       RelationshipType = "hasOptionalComponent"
	RelationshipTypeHasOptionalDependency      RelationshipType = "hasOptionalDependency"
	RelationshipTypeHasOutput                  RelationshipType = "hasOutput"
	RelationshipTypeHasPrerequisite            RelationshipType = "hasPrerequisite"
	RelationshipTypeHasProvidedDependency      RelationshipType = "hasProvidedDependency"
	RelationshipTypeHasRequirement             RelationshipType = "hasRequirement"
	RelationshipTypeHasSpecification           RelationshipType = "hasSpecification"
	RelationshipTypeHasStaticLink              RelationshipType = "hasStaticLink"
	RelationshipTypeHasTest                    RelationshipType = "hasTest"
	RelationshipTypeHasTestCase                RelationshipType = "hasTestCase"
	RelationshipTypeHasVariant                 RelationshipType = "hasVariant"
	RelationshipTypeInvokedBy                  RelationshipType = "invokedBy"
	RelationshipTypeModifiedBy                 RelationshipType = "modifiedBy"
	RelationshipTypeOther                      RelationshipType = "other"
	RelationshipTypePackagedBy                 RelationshipType = "packagedBy"
	RelationshipTypePatchedBy                  RelationshipType = "patchedBy"
	RelationshipTypePublishedBy                RelationshipType = "publishedBy"
	RelationshipTypeReportedBy                 RelationshipType = "reportedBy"
	RelationshipTypeRepublishedBy              RelationshipType = "republishedBy"
	RelationshipTypeSerializedInArtifact       RelationshipType = "serializedInArtifact"
	RelationshipTypeTestedOn                   RelationshipType = "testedOn"
	RelationshipTypeTrainedOn                  RelationshipType = "trainedOn"
	RelationshipTypeUnderInvestigationFor      RelationshipType = "underInvestigationFor"
	RelationshipTypeUsesTool                   RelationshipType = "usesTool"
)

// SupportType vocabulary defines valid support types
type SupportType string

var SupportTypes = types.Vocabulary[SupportType]{
	SupportTypeDeployed,
	SupportTypeDevelopment,
	SupportTypeEndOfSupport,
	SupportTypeLimitedSupport,
	SupportTypeNoAssertion,
	SupportTypeNoSupport,
	SupportTypeSupport,
}

const (
	// SupportTypeDeployed indicates software is deployed and actively used
	SupportTypeDeployed SupportType = "deployed"
	// SupportTypeDevelopment indicates active development phase
	SupportTypeDevelopment SupportType = "development"
	// SupportTypeEndOfSupport indicates defined end of support period
	SupportTypeEndOfSupport SupportType = "endOfSupport"
	// SupportTypeLimitedSupport indicates released with limited supplier support
	SupportTypeLimitedSupport SupportType = "limitedSupport"
	// SupportTypeNoAssertion indicates no assertion about support type
	SupportTypeNoAssertion SupportType = "noAssertion"
	// SupportTypeNoSupport indicates no supplier support
	SupportTypeNoSupport SupportType = "noSupport"
	// SupportTypeSupport indicates released with active supplier support
	SupportTypeSupport SupportType = "support"
)
