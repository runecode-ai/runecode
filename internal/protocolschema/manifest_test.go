package protocolschema

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v6"
)

const (
	bundleID            = "runecode.protocol.v0"
	runtimeSchemaPrefix = bundleID + "."
	manifestMetaPath    = "meta/manifest.schema.json"
	registryMetaPath    = "meta/registry.schema.json"
)

var (
	allowedDataClasses = map[string]struct{}{
		"public":    {},
		"sensitive": {},
		"secret":    {},
	}
	placeholderSchemaIDs = map[string]struct{}{
		"runecode.protocol.v0.ApprovalRequest":  {},
		"runecode.protocol.v0.ApprovalDecision": {},
		"runecode.protocol.v0.PolicyDecision":   {},
		"runecode.protocol.v0.Error":            {},
	}
)

type manifestFile struct {
	BundleID                   string                `json:"bundle_id"`
	BundleVersion              string                `json:"bundle_version"`
	JSONSchemaDraft            string                `json:"json_schema_draft"`
	RuntimeSchemaPrefix        string                `json:"runtime_schema_prefix"`
	Canonicalization           string                `json:"canonicalization"`
	TopLevelObjectRequirements topLevelRequirements  `json:"top_level_object_requirements"`
	SchemaFiles                []schemaManifestEntry `json:"schema_files"`
	Registries                 []registryManifest    `json:"registries"`
}

type topLevelRequirements struct {
	RequireSchemaID      bool   `json:"require_schema_id"`
	RequireSchemaVersion bool   `json:"require_schema_version"`
	UnknownSchemaPosture string `json:"unknown_schema_posture"`
}

type schemaManifestEntry struct {
	Path          string `json:"path"`
	SchemaID      string `json:"schema_id"`
	SchemaVersion string `json:"schema_version"`
	Owner         string `json:"owner"`
	Status        string `json:"status"`
	Note          string `json:"note"`
}

type registryManifest struct {
	Path               string `json:"path"`
	Name               string `json:"name"`
	Namespace          string `json:"namespace"`
	DocumentationOwner string `json:"documentation_owner"`
	Status             string `json:"status"`
}

type registryFile struct {
	RegistryName       string         `json:"registry_name"`
	Namespace          string         `json:"namespace"`
	DocumentationOwner string         `json:"documentation_owner"`
	Status             string         `json:"status"`
	Description        string         `json:"description"`
	Codes              []registryCode `json:"codes"`
}

type registryCode struct {
	Code    string `json:"code"`
	Summary string `json:"summary"`
}

type compiledBundle struct {
	Compiler   *jsonschema.Compiler
	SchemaDocs map[string]map[string]any
}

type resolvedSchemaRef struct {
	FilePath string
	Location string
	Node     map[string]any
}

func TestSchemaManifestMatchesSchemas(t *testing.T) {
	manifest := loadManifest(t)

	if manifest.BundleID != bundleID {
		t.Fatalf("bundle_id = %q, want %q", manifest.BundleID, bundleID)
	}

	if manifest.BundleVersion == "" {
		t.Fatal("bundle_version must be non-empty")
	}

	if manifest.JSONSchemaDraft != "2020-12" {
		t.Fatalf("json_schema_draft = %q, want 2020-12", manifest.JSONSchemaDraft)
	}

	if manifest.RuntimeSchemaPrefix != runtimeSchemaPrefix {
		t.Fatalf("runtime_schema_prefix = %q, want %q", manifest.RuntimeSchemaPrefix, runtimeSchemaPrefix)
	}

	if manifest.Canonicalization != "RFC8785-JCS" {
		t.Fatalf("canonicalization = %q, want RFC8785-JCS", manifest.Canonicalization)
	}

	if !manifest.TopLevelObjectRequirements.RequireSchemaID {
		t.Fatal("top-level objects must require schema_id")
	}

	if !manifest.TopLevelObjectRequirements.RequireSchemaVersion {
		t.Fatal("top-level objects must require schema_version")
	}

	if manifest.TopLevelObjectRequirements.UnknownSchemaPosture != "fail_closed" {
		t.Fatalf("unknown_schema_posture = %q, want fail_closed", manifest.TopLevelObjectRequirements.UnknownSchemaPosture)
	}

	assertManifestFileSet(t, manifest.SchemaFiles, "objects", ".schema.json")
	assertManifestRegistryFileSet(t, manifest.Registries)

	seenIDs := map[string]string{}
	for _, entry := range manifest.SchemaFiles {
		t.Run(entry.Path, func(t *testing.T) {
			if previous, ok := seenIDs[entry.SchemaID]; ok {
				t.Fatalf("duplicate schema_id %q in %q and %q", entry.SchemaID, previous, entry.Path)
			}
			seenIDs[entry.SchemaID] = entry.Path

			if !strings.HasPrefix(entry.SchemaID, manifest.RuntimeSchemaPrefix) {
				t.Fatalf("schema_id %q does not use runtime prefix %q", entry.SchemaID, manifest.RuntimeSchemaPrefix)
			}

			if entry.SchemaVersion == "" {
				t.Fatalf("schema_version for %q must be non-empty", entry.Path)
			}

			if entry.Owner != "protocol" {
				t.Fatalf("owner for %q = %q, want protocol", entry.Path, entry.Owner)
			}

			if entry.Status != "mvp" && entry.Status != "reserved" {
				t.Fatalf("status for %q = %q, want mvp or reserved", entry.Path, entry.Status)
			}

			if requiresPlaceholderNote(entry.SchemaID) && strings.TrimSpace(entry.Note) == "" {
				t.Fatalf("schema %q must carry a manifest note explaining its placeholder scope", entry.SchemaID)
			}

			schema := loadJSONMap(t, schemaPath(t, entry.Path))
			if got := stringValue(t, schema, "$schema"); got != "https://json-schema.org/draft/2020-12/schema" {
				t.Fatalf("$schema for %q = %q, want draft 2020-12", entry.Path, got)
			}

			if got := stringValue(t, schema, "$id"); got == "" {
				t.Fatalf("$id for %q must be non-empty", entry.Path)
			}

			if got := stringValue(t, schema, "type"); got != "object" {
				t.Fatalf("type for %q = %q, want object", entry.Path, got)
			}

			if boolValue(t, schema, "additionalProperties") {
				t.Fatalf("additionalProperties for %q must be false", entry.Path)
			}

			if !hasNumber(schema, "maxProperties") {
				t.Fatalf("schema %q must declare maxProperties", entry.Path)
			}

			required := stringSliceValue(t, schema, "required")
			assertContains(t, required, "schema_id")
			assertContains(t, required, "schema_version")

			properties := objectValue(t, schema, "properties")
			assertConst(t, properties, "schema_id", entry.SchemaID)
			assertConst(t, properties, "schema_version", entry.SchemaVersion)
		})
	}

	assertReservedStatus(t, manifest, "runecode.protocol.v0.WorkflowDefinition")
	assertReservedStatus(t, manifest, "runecode.protocol.v0.ProcessDefinition")
}

func TestManifestAndRegistryDocumentsValidateAgainstMetaSchemas(t *testing.T) {
	manifest := loadManifest(t)
	compiler := newMetaCompiler(t)

	manifestSchema := mustCompileMetaSchema(t, compiler, metaPath(t, manifestMetaPath))
	if err := manifestSchema.Validate(loadJSONMap(t, schemaPath(t, "manifest.json"))); err != nil {
		t.Fatalf("manifest.json failed meta-schema validation: %v", err)
	}

	registrySchema := mustCompileMetaSchema(t, compiler, metaPath(t, registryMetaPath))
	for _, entry := range manifest.Registries {
		t.Run(entry.Path, func(t *testing.T) {
			if err := registrySchema.Validate(loadJSONMap(t, schemaPath(t, entry.Path))); err != nil {
				t.Fatalf("%s failed registry meta-schema validation: %v", entry.Path, err)
			}
		})
	}
}

func TestSchemasCompileAgainstDraft202012(t *testing.T) {
	bundle := newCompiledBundle(t, loadManifest(t))

	for filePath, schemaDoc := range bundle.SchemaDocs {
		t.Run(filePath, func(t *testing.T) {
			schemaID := stringValue(t, schemaDoc, "$id")
			if _, err := bundle.Compiler.Compile(schemaID); err != nil {
				t.Fatalf("Compile(%q) returned error: %v", filePath, err)
			}
		})
	}
}

func TestSchemaPropertiesHaveClassificationBoundsAndDescriptions(t *testing.T) {
	bundle := newCompiledBundle(t, loadManifest(t))

	for filePath, schemaDoc := range bundle.SchemaDocs {
		t.Run(filePath, func(t *testing.T) {
			assertSchemaNodeInvariants(t, filePath, schemaDoc, false)
			assertReferencedDefinitions(t, filePath, schemaDoc, bundle.SchemaDocs, map[string]struct{}{})
		})
	}
}

func TestRegistryNamespacesAreSeparate(t *testing.T) {
	manifest := loadManifest(t)
	seenNames := map[string]struct{}{}
	seenNamespaces := map[string]struct{}{}
	codesByRegistry := map[string]map[string]struct{}{}
	registryNames := make([]string, 0, len(manifest.Registries))

	for _, entry := range manifest.Registries {
		t.Run(entry.Path, func(t *testing.T) {
			if _, ok := seenNames[entry.Name]; ok {
				t.Fatalf("duplicate registry name %q", entry.Name)
			}
			seenNames[entry.Name] = struct{}{}
			registryNames = append(registryNames, entry.Name)

			if _, ok := seenNamespaces[entry.Namespace]; ok {
				t.Fatalf("duplicate registry namespace %q", entry.Namespace)
			}
			seenNamespaces[entry.Namespace] = struct{}{}

			if entry.DocumentationOwner != "protocol" {
				t.Fatalf("documentation_owner for %q = %q, want protocol", entry.Path, entry.DocumentationOwner)
			}

			if entry.Status != "mvp" {
				t.Fatalf("status for %q = %q, want mvp", entry.Path, entry.Status)
			}

			registry := loadRegistry(t, schemaPath(t, entry.Path))
			if registry.RegistryName != entry.Name {
				t.Fatalf("registry_name for %q = %q, want %q", entry.Path, registry.RegistryName, entry.Name)
			}

			if registry.Namespace != entry.Namespace {
				t.Fatalf("namespace for %q = %q, want %q", entry.Path, registry.Namespace, entry.Namespace)
			}

			if registry.DocumentationOwner != entry.DocumentationOwner {
				t.Fatalf("documentation_owner for %q = %q, want %q", entry.Path, registry.DocumentationOwner, entry.DocumentationOwner)
			}

			if registry.Status != entry.Status {
				t.Fatalf("status for %q = %q, want %q", entry.Path, registry.Status, entry.Status)
			}

			if strings.TrimSpace(registry.Description) == "" {
				t.Fatalf("registry %q must have a non-empty description", entry.Name)
			}

			seenCodes := map[string]struct{}{}
			codesByRegistry[entry.Name] = map[string]struct{}{}
			for _, code := range registry.Codes {
				if code.Code == "" {
					t.Fatalf("registry %q has empty code", entry.Name)
				}

				if _, ok := seenCodes[code.Code]; ok {
					t.Fatalf("registry %q reuses code %q", entry.Name, code.Code)
				}
				seenCodes[code.Code] = struct{}{}
				codesByRegistry[entry.Name][code.Code] = struct{}{}

				if strings.TrimSpace(code.Summary) == "" {
					t.Fatalf("registry %q code %q must have a non-empty summary", entry.Name, code.Code)
				}
			}
		})
	}

	sort.Strings(registryNames)
	for i := 0; i < len(registryNames); i++ {
		for j := i + 1; j < len(registryNames); j++ {
			assertNoCodeOverlap(t, codesByRegistry, registryNames[i], registryNames[j])
		}
	}

	errorRegistry := loadRegistry(t, schemaPath(t, "registries/error.code.registry.json"))
	assertRegistryCode(t, errorRegistry, "unknown_schema_id")
	assertRegistryCode(t, errorRegistry, "unsupported_schema_version")
	assertRegistryCode(t, errorRegistry, "unsupported_hash_algorithm")
}

func TestTaskTwoSchemaRequirements(t *testing.T) {
	bundle := newCompiledBundle(t, loadManifest(t))

	t.Run("digest schema pins sha256", func(t *testing.T) {
		schema := loadJSONMap(t, schemaPath(t, "objects/Digest.schema.json"))
		required := stringSliceValue(t, schema, "required")
		assertContains(t, required, "hash_alg")
		assertContains(t, required, "hash")

		properties := objectValue(t, schema, "properties")
		assertConst(t, properties, "hash_alg", "sha256")

		digestValue := objectValue(t, objectValue(t, schema, "$defs"), "digestValue")
		digestRequired := stringSliceValue(t, digestValue, "required")
		assertContains(t, digestRequired, "hash_alg")
		assertContains(t, digestRequired, "hash")
	})

	t.Run("signed envelope constrains payload and algorithms", func(t *testing.T) {
		schema := loadJSONMap(t, schemaPath(t, "objects/SignedObjectEnvelope.schema.json"))
		required := stringSliceValue(t, schema, "required")
		assertContains(t, required, "payload")
		assertContains(t, required, "signature_input")
		assertContains(t, required, "signature")

		properties := objectValue(t, schema, "properties")
		assertConst(t, properties, "signature_input", "rfc8785_jcs_detached_payload")

		payload := objectValue(t, properties, "payload")
		if got := stringValue(t, payload, "type"); got != "object" {
			t.Fatalf("payload type = %q, want object", got)
		}
		payloadRequired := stringSliceValue(t, payload, "required")
		assertContains(t, payloadRequired, "schema_id")
		assertContains(t, payloadRequired, "schema_version")
		if got := stringValue(t, payload, "x-data-class"); got != "secret" {
			t.Fatalf("payload x-data-class = %q, want secret", got)
		}

		signatureBlock := objectValue(t, objectValue(t, schema, "$defs"), "signatureBlock")
		signatureRequired := stringSliceValue(t, signatureBlock, "required")
		assertContains(t, signatureRequired, "alg")
		assertContains(t, signatureRequired, "key_id")
		assertContains(t, signatureRequired, "signature")

		alg := objectValue(t, objectValue(t, signatureBlock, "properties"), "alg")
		assertContains(t, stringSliceValue(t, alg, "enum"), "ed25519")
	})

	t.Run("manifests require explicit signed inputs", func(t *testing.T) {
		for _, schemaFile := range []string{
			"objects/RoleManifest.schema.json",
			"objects/CapabilityManifest.schema.json",
		} {
			t.Run(schemaFile, func(t *testing.T) {
				schema := loadJSONMap(t, schemaPath(t, schemaFile))
				required := stringSliceValue(t, schema, "required")
				assertContains(t, required, "approval_profile")
				assertContains(t, required, "capability_opt_ins")
				assertContains(t, required, "allowlist_refs")
				assertContains(t, required, "signatures")
			})
		}

		capabilitySchema := loadJSONMap(t, schemaPath(t, "objects/CapabilityManifest.schema.json"))
		properties := objectValue(t, capabilitySchema, "properties")
		manifestScope := objectValue(t, properties, "manifest_scope")
		enumValues := stringSliceValue(t, manifestScope, "enum")
		assertContains(t, enumValues, "run")
		assertContains(t, enumValues, "stage")
	})

	t.Run("principal identity constrains role_kind by actor kind", func(t *testing.T) {
		schema := mustCompileObjectSchema(t, bundle, "objects/PrincipalIdentity.schema.json")

		tests := []struct {
			name    string
			value   map[string]any
			wantErr bool
		}{
			{
				name: "role instance requires role kind",
				value: map[string]any{
					"schema_id":      "runecode.protocol.v0.PrincipalIdentity",
					"schema_version": "0.1.0",
					"actor_kind":     "role_instance",
					"principal_id":   "role-123",
					"role_kind":      "gateway",
				},
			},
			{
				name: "role instance without role kind fails",
				value: map[string]any{
					"schema_id":      "runecode.protocol.v0.PrincipalIdentity",
					"schema_version": "0.1.0",
					"actor_kind":     "role_instance",
					"principal_id":   "role-123",
				},
				wantErr: true,
			},
			{
				name: "daemon may include role kind",
				value: map[string]any{
					"schema_id":      "runecode.protocol.v0.PrincipalIdentity",
					"schema_version": "0.1.0",
					"actor_kind":     "daemon",
					"principal_id":   "secretsd",
					"role_kind":      "auth",
				},
			},
			{
				name: "user may not include role kind",
				value: map[string]any{
					"schema_id":      "runecode.protocol.v0.PrincipalIdentity",
					"schema_version": "0.1.0",
					"actor_kind":     "user",
					"principal_id":   "alice",
					"role_kind":      "gateway",
				},
				wantErr: true,
			},
			{
				name: "local client may not include role kind",
				value: map[string]any{
					"schema_id":      "runecode.protocol.v0.PrincipalIdentity",
					"schema_version": "0.1.0",
					"actor_kind":     "local_client",
					"principal_id":   "cli-session",
					"role_kind":      "workspace",
				},
				wantErr: true,
			},
		}

		for _, testCase := range tests {
			t.Run(testCase.name, func(t *testing.T) {
				err := schema.Validate(testCase.value)
				if testCase.wantErr && err == nil {
					t.Fatal("Validate returned nil error, want failure")
				}
				if !testCase.wantErr && err != nil {
					t.Fatalf("Validate returned error: %v", err)
				}
			})
		}
	})
}

func loadManifest(t *testing.T) manifestFile {
	t.Helper()

	var manifest manifestFile
	loadJSON(t, schemaPath(t, "manifest.json"), &manifest)
	return manifest
}

func loadRegistry(t *testing.T, filePath string) registryFile {
	t.Helper()

	var registry registryFile
	loadJSON(t, filePath, &registry)
	return registry
}

func newCompiledBundle(t *testing.T, manifest manifestFile) compiledBundle {
	t.Helper()

	compiler := jsonschema.NewCompiler()
	schemaDocs := make(map[string]map[string]any, len(manifest.SchemaFiles))

	for _, entry := range manifest.SchemaFiles {
		schemaDoc := loadJSONMap(t, schemaPath(t, entry.Path))
		schemaID := stringValue(t, schemaDoc, "$id")
		if err := compiler.AddResource(schemaID, schemaDoc); err != nil {
			t.Fatalf("AddResource(%q) returned error: %v", schemaID, err)
		}
		schemaDocs[entry.Path] = schemaDoc
	}

	return compiledBundle{Compiler: compiler, SchemaDocs: schemaDocs}
}

func newMetaCompiler(t *testing.T) *jsonschema.Compiler {
	t.Helper()

	compiler := jsonschema.NewCompiler()
	for _, metaFile := range []string{manifestMetaPath, registryMetaPath} {
		doc := loadJSONMap(t, metaPath(t, metaFile))
		docID := stringValue(t, doc, "$id")
		if err := compiler.AddResource(docID, doc); err != nil {
			t.Fatalf("AddResource(%q) returned error: %v", docID, err)
		}
	}

	return compiler
}

func mustCompileMetaSchema(t *testing.T, compiler *jsonschema.Compiler, filePath string) *jsonschema.Schema {
	t.Helper()

	doc := loadJSONMap(t, filePath)
	docID := stringValue(t, doc, "$id")
	schema, err := compiler.Compile(docID)
	if err != nil {
		t.Fatalf("Compile(%q) returned error: %v", filePath, err)
	}
	return schema
}

func mustCompileObjectSchema(t *testing.T, bundle compiledBundle, filePath string) *jsonschema.Schema {
	t.Helper()

	doc, ok := bundle.SchemaDocs[filePath]
	if !ok {
		t.Fatalf("schema document %q not found", filePath)
	}
	objID := stringValue(t, doc, "$id")
	schema, err := bundle.Compiler.Compile(objID)
	if err != nil {
		t.Fatalf("Compile(%q) returned error: %v", filePath, err)
	}
	return schema
}

func assertSchemaNodeInvariants(t *testing.T, location string, node map[string]any, requireClassification bool) {
	t.Helper()

	if requireClassification {
		description := strings.TrimSpace(stringValue(t, node, "description"))
		if description == "" {
			t.Fatalf("%s must have a non-empty description", location)
		}

		dataClass := stringValue(t, node, "x-data-class")
		if _, ok := allowedDataClasses[dataClass]; !ok {
			t.Fatalf("%s uses unsupported x-data-class %q", location, dataClass)
		}
	}

	if schemaType, ok := optionalStringValue(node, "type"); ok {
		switch schemaType {
		case "object":
			if !hasNumber(node, "maxProperties") {
				t.Fatalf("%s must declare maxProperties", location)
			}
		case "array":
			if !hasNumber(node, "maxItems") {
				t.Fatalf("%s must declare maxItems", location)
			}
		case "string":
			if !hasNumber(node, "maxLength") && !hasKey(node, "const") && !hasKey(node, "enum") {
				t.Fatalf("%s must declare maxLength or constrain values with const/enum", location)
			}
		}
	}

	if properties, ok := optionalObjectValue(node, "properties"); ok {
		for _, key := range sortedKeys(properties) {
			child := objectFromAny(t, location+"."+key, properties[key])
			assertSchemaNodeInvariants(t, location+"."+key, child, true)
		}
	}

	if defs, ok := optionalObjectValue(node, "$defs"); ok {
		for _, key := range sortedKeys(defs) {
			child := objectFromAny(t, location+".$defs."+key, defs[key])
			if strings.TrimSpace(stringValue(t, child, "description")) == "" {
				t.Fatalf("%s.$defs.%s must have a non-empty description", location, key)
			}
			assertSchemaNodeInvariants(t, location+".$defs."+key, child, false)
		}
	}

	if items, ok := optionalObjectValue(node, "items"); ok {
		assertSchemaNodeInvariants(t, location+"[]", items, false)
	}
}

func assertReferencedDefinitions(t *testing.T, currentFile string, node map[string]any, schemaDocs map[string]map[string]any, seen map[string]struct{}) {
	t.Helper()

	if ref, ok := optionalStringValue(node, "$ref"); ok {
		resolved := resolveSchemaRef(t, currentFile, ref, schemaDocs)
		if _, ok := seen[resolved.Location]; !ok {
			seen[resolved.Location] = struct{}{}
			assertSchemaNodeInvariants(t, resolved.Location, resolved.Node, false)
			assertReferencedDefinitions(t, resolved.FilePath, resolved.Node, schemaDocs, seen)
		}
	}

	if properties, ok := optionalObjectValue(node, "properties"); ok {
		for _, key := range sortedKeys(properties) {
			child := objectFromAny(t, currentFile+"."+key, properties[key])
			assertReferencedDefinitions(t, currentFile, child, schemaDocs, seen)
		}
	}

	if defs, ok := optionalObjectValue(node, "$defs"); ok {
		for _, key := range sortedKeys(defs) {
			child := objectFromAny(t, currentFile+".$defs."+key, defs[key])
			assertReferencedDefinitions(t, currentFile, child, schemaDocs, seen)
		}
	}

	if items, ok := optionalObjectValue(node, "items"); ok {
		assertReferencedDefinitions(t, currentFile, items, schemaDocs, seen)
	}
}

func resolveSchemaRef(t *testing.T, currentFile string, ref string, schemaDocs map[string]map[string]any) resolvedSchemaRef {
	t.Helper()

	refPath := currentFile
	fragment := ""
	if hashIndex := strings.IndexByte(ref, '#'); hashIndex >= 0 {
		if hashIndex > 0 {
			refPath = path.Clean(path.Join(path.Dir(currentFile), ref[:hashIndex]))
		}
		fragment = ref[hashIndex+1:]
	} else {
		refPath = path.Clean(path.Join(path.Dir(currentFile), ref))
	}

	doc, ok := schemaDocs[refPath]
	if !ok {
		t.Fatalf("reference %q from %q resolved to unknown schema file %q", ref, currentFile, refPath)
	}

	resolvedNode := any(doc)
	location := refPath
	if fragment != "" {
		resolvedNode = resolveJSONPointer(t, doc, fragment)
		location = fmt.Sprintf("%s#%s", refPath, fragment)
	}

	objectNode, ok := resolvedNode.(map[string]any)
	if !ok {
		t.Fatalf("reference %q from %q resolved to %T, want map[string]any", ref, currentFile, resolvedNode)
	}

	return resolvedSchemaRef{FilePath: refPath, Location: location, Node: objectNode}
}

func resolveJSONPointer(t *testing.T, value any, pointer string) any {
	t.Helper()

	if pointer == "" {
		return value
	}
	if !strings.HasPrefix(pointer, "/") {
		t.Fatalf("json pointer %q must begin with '/'", pointer)
	}

	current := value
	for _, rawToken := range strings.Split(pointer[1:], "/") {
		token := strings.ReplaceAll(strings.ReplaceAll(rawToken, "~1", "/"), "~0", "~")
		switch typed := current.(type) {
		case map[string]any:
			next, ok := typed[token]
			if !ok {
				t.Fatalf("json pointer %q segment %q not found", pointer, token)
			}
			current = next
		case []any:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(typed) {
				t.Fatalf("json pointer %q segment %q is not a valid array index", pointer, token)
			}
			current = typed[index]
		default:
			t.Fatalf("json pointer %q cannot descend into %T", pointer, current)
		}
	}

	return current
}

func assertManifestFileSet(t *testing.T, entries []schemaManifestEntry, dir string, suffix string) {
	t.Helper()

	manifestPaths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Path, dir+"/") {
			t.Fatalf("manifest path %q must stay under %s/", entry.Path, dir)
		}
		if !strings.HasSuffix(entry.Path, suffix) {
			t.Fatalf("manifest path %q must end with %q", entry.Path, suffix)
		}
		_ = schemaPath(t, entry.Path)
		manifestPaths = append(manifestPaths, entry.Path)
	}

	actualPaths := listedFiles(t, schemaRoot(), dir, suffix)
	assertSameStringSet(t, manifestPaths, actualPaths)
}

func assertManifestRegistryFileSet(t *testing.T, entries []registryManifest) {
	t.Helper()

	manifestPaths := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Path, "registries/") {
			t.Fatalf("registry path %q must stay under registries/", entry.Path)
		}
		if !strings.HasSuffix(entry.Path, ".registry.json") {
			t.Fatalf("registry path %q must end with .registry.json", entry.Path)
		}
		_ = schemaPath(t, entry.Path)
		manifestPaths = append(manifestPaths, entry.Path)
	}

	actualPaths := listedFiles(t, schemaRoot(), "registries", ".registry.json")
	assertSameStringSet(t, manifestPaths, actualPaths)
}

func listedFiles(t *testing.T, root string, dir string, suffix string) []string {
	t.Helper()

	entries, err := os.ReadDir(filepath.Join(root, dir))
	if err != nil {
		t.Fatalf("ReadDir(%q) returned error: %v", filepath.Join(root, dir), err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), suffix) {
			files = append(files, path.Join(dir, entry.Name()))
		}
	}

	return files
}

func assertSameStringSet(t *testing.T, got []string, want []string) {
	t.Helper()

	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("set size mismatch: got %v, want %v", got, want)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("set mismatch: got %v, want %v", got, want)
		}
	}
}

func loadJSON(t *testing.T, filePath string, target any) {
	t.Helper()

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("ReadFile(%q) returned error: %v", filePath, err)
	}

	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("Unmarshal(%q) returned error: %v", filePath, err)
	}
}

func loadJSONMap(t *testing.T, filePath string) map[string]any {
	t.Helper()

	var value map[string]any
	loadJSON(t, filePath, &value)
	return value
}

func schemaRoot() string {
	return filepath.Join("..", "..", "protocol", "schemas")
}

func metaRoot() string {
	return filepath.Join(schemaRoot(), "meta")
}

func schemaPath(t *testing.T, rel string) string {
	t.Helper()

	return rootedSchemaPath(t, schemaRoot(), rel, "protocol/schemas")
}

func metaPath(t *testing.T, rel string) string {
	t.Helper()

	return rootedSchemaPath(t, schemaRoot(), rel, "protocol/schemas")
}

func rootedSchemaPath(t *testing.T, root string, rel string, label string) string {
	t.Helper()

	if rel == "" {
		t.Fatalf("%s path must be non-empty", label)
	}

	if filepath.IsAbs(rel) || path.IsAbs(rel) {
		t.Fatalf("%s path %q must be relative", label, rel)
	}

	cleaned := path.Clean(rel)
	if cleaned != rel {
		t.Fatalf("%s path %q must already be clean; got %q", label, rel, cleaned)
	}

	if cleaned == ".." || strings.HasPrefix(cleaned, "../") {
		t.Fatalf("%s path %q escapes %s", label, rel, label)
	}

	absPath := filepath.Join(root, filepath.FromSlash(cleaned))
	relToRoot, err := filepath.Rel(root, absPath)
	if err != nil {
		t.Fatalf("Rel(%q) returned error: %v", rel, err)
	}

	if relToRoot == ".." || strings.HasPrefix(relToRoot, ".."+string(filepath.Separator)) {
		t.Fatalf("%s path %q escapes %s", label, rel, label)
	}

	return absPath
}

func stringValue(t *testing.T, object map[string]any, key string) string {
	t.Helper()

	value, ok := object[key]
	if !ok {
		t.Fatalf("missing key %q", key)
	}

	stringValue, ok := value.(string)
	if !ok {
		t.Fatalf("key %q has type %T, want string", key, value)
	}

	return stringValue
}

func optionalStringValue(object map[string]any, key string) (string, bool) {
	value, ok := object[key]
	if !ok {
		return "", false
	}

	stringValue, ok := value.(string)
	return stringValue, ok
}

func boolValue(t *testing.T, object map[string]any, key string) bool {
	t.Helper()

	value, ok := object[key]
	if !ok {
		t.Fatalf("missing key %q", key)
	}

	boolValue, ok := value.(bool)
	if !ok {
		t.Fatalf("key %q has type %T, want bool", key, value)
	}

	return boolValue
}

func stringSliceValue(t *testing.T, object map[string]any, key string) []string {
	t.Helper()

	value, ok := object[key]
	if !ok {
		t.Fatalf("missing key %q", key)
	}

	items, ok := value.([]any)
	if !ok {
		t.Fatalf("key %q has type %T, want []any", key, value)
	}

	result := make([]string, 0, len(items))
	for _, item := range items {
		stringItem, ok := item.(string)
		if !ok {
			t.Fatalf("key %q has non-string item type %T", key, item)
		}
		result = append(result, stringItem)
	}

	return result
}

func objectValue(t *testing.T, object map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := object[key]
	if !ok {
		t.Fatalf("missing key %q", key)
	}

	child, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("key %q has type %T, want map[string]any", key, value)
	}

	return child
}

func optionalObjectValue(object map[string]any, key string) (map[string]any, bool) {
	value, ok := object[key]
	if !ok {
		return nil, false
	}

	child, ok := value.(map[string]any)
	return child, ok
}

func objectFromAny(t *testing.T, location string, value any) map[string]any {
	t.Helper()

	child, ok := value.(map[string]any)
	if !ok {
		t.Fatalf("%s has type %T, want map[string]any", location, value)
	}

	return child
}

func hasKey(object map[string]any, key string) bool {
	_, ok := object[key]
	return ok
}

func hasNumber(object map[string]any, key string) bool {
	value, ok := object[key]
	if !ok {
		return false
	}
	_, ok = value.(float64)
	return ok
}

func sortedKeys(object map[string]any) []string {
	keys := make([]string, 0, len(object))
	for key := range object {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func assertContains(t *testing.T, values []string, want string) {
	t.Helper()

	for _, value := range values {
		if value == want {
			return
		}
	}

	t.Fatalf("%q not found in %v", want, values)
}

func assertConst(t *testing.T, properties map[string]any, key string, want string) {
	t.Helper()

	property := objectValue(t, properties, key)
	if got := stringValue(t, property, "const"); got != want {
		t.Fatalf("const for property %q = %q, want %q", key, got, want)
	}
}

func assertReservedStatus(t *testing.T, manifest manifestFile, schemaID string) {
	t.Helper()

	for _, entry := range manifest.SchemaFiles {
		if entry.SchemaID == schemaID {
			if entry.Status != "reserved" {
				t.Fatalf("status for %q = %q, want reserved", schemaID, entry.Status)
			}
			return
		}
	}

	t.Fatalf("schema_id %q not found in manifest", schemaID)
}

func assertRegistryCode(t *testing.T, registry registryFile, want string) {
	t.Helper()

	for _, code := range registry.Codes {
		if code.Code == want {
			return
		}
	}

	t.Fatalf("registry %q missing code %q", registry.RegistryName, want)
}

func assertNoCodeOverlap(t *testing.T, codesByRegistry map[string]map[string]struct{}, left string, right string) {
	t.Helper()

	for code := range codesByRegistry[left] {
		if _, ok := codesByRegistry[right][code]; ok {
			t.Fatalf("registry code %q must not appear in both %q and %q", code, left, right)
		}
	}
}

func requiresPlaceholderNote(schemaID string) bool {
	_, ok := placeholderSchemaIDs[schemaID]
	return ok
}
