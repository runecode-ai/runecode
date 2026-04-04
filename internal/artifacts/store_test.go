package artifacts

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPutGetHeadAndList(t *testing.T) {
	store := newTestStore(t)
	ref, err := store.Put(PutRequest{
		Payload:               []byte("hello"),
		ContentType:           "text/plain",
		DataClass:             DataClassSpecText,
		ProvenanceReceiptHash: testDigest("1"),
		CreatedByRole:         "workspace",
		RunID:                 "run-a",
		StepID:                "step-a",
	})
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	record, err := store.Head(ref.Digest)
	if err != nil {
		t.Fatalf("Head returned error: %v", err)
	}
	if record.Reference.DataClass != DataClassSpecText {
		t.Fatalf("Head data class = %q, want %q", record.Reference.DataClass, DataClassSpecText)
	}
	r, err := store.Get(ref.Digest)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	b, readErr := ioReadAllAndClose(r)
	if readErr != nil {
		t.Fatalf("read payload error: %v", readErr)
	}
	if string(b) != "hello" {
		t.Fatalf("Get payload = %q, want hello", string(b))
	}
	if len(store.List()) != 1 {
		t.Fatalf("List count = %d, want 1", len(store.List()))
	}
}

func TestCanonicalJSONDigestDeterministic(t *testing.T) {
	store := newTestStore(t)
	ref1, err := store.Put(PutRequest{
		Payload:               []byte(`{"b":2,"a":1}`),
		ContentType:           "application/json",
		DataClass:             DataClassSpecText,
		ProvenanceReceiptHash: testDigest("2"),
		CreatedByRole:         "workspace",
	})
	if err != nil {
		t.Fatalf("first Put returned error: %v", err)
	}
	ref2, err := store.Put(PutRequest{
		Payload:               []byte(`{"a":1,"b":2}`),
		ContentType:           "application/json",
		DataClass:             DataClassSpecText,
		ProvenanceReceiptHash: testDigest("2"),
		CreatedByRole:         "workspace",
	})
	if err != nil {
		t.Fatalf("second Put returned error: %v", err)
	}
	if ref1.Digest != ref2.Digest {
		t.Fatalf("digests differ: %s vs %s", ref1.Digest, ref2.Digest)
	}
}

func TestDataClassMutationDenied(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Put(PutRequest{
		Payload:               []byte("same-bytes"),
		ContentType:           "text/plain",
		DataClass:             DataClassSpecText,
		ProvenanceReceiptHash: testDigest("3"),
		CreatedByRole:         "workspace",
	})
	if err != nil {
		t.Fatalf("first Put returned error: %v", err)
	}
	_, err = store.Put(PutRequest{
		Payload:               []byte("same-bytes"),
		ContentType:           "text/plain",
		DataClass:             DataClassDiffs,
		ProvenanceReceiptHash: testDigest("3"),
		CreatedByRole:         "workspace",
	})
	if err != ErrDataClassMutationDenied {
		t.Fatalf("Put error = %v, want %v", err, ErrDataClassMutationDenied)
	}
}

func TestFlowChecksFailClosedAndEgressRules(t *testing.T) {
	store := newTestStore(t)
	ref, err := store.Put(PutRequest{
		Payload:               []byte("excerpt"),
		ContentType:           "text/plain",
		DataClass:             DataClassUnapprovedFileExcerpts,
		ProvenanceReceiptHash: testDigest("4"),
		CreatedByRole:         "workspace",
	})
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	err = store.CheckFlow(FlowCheckRequest{ProducerRole: "workspace", ConsumerRole: "model_gateway", DataClass: DataClassUnapprovedFileExcerpts, Digest: ref.Digest, IsEgress: true})
	if err != ErrUnapprovedEgressDenied {
		t.Fatalf("CheckFlow error = %v, want %v", err, ErrUnapprovedEgressDenied)
	}

	err = store.CheckFlow(FlowCheckRequest{ProducerRole: "workspace", ConsumerRole: "model_gateway", DataClass: DataClassApprovedFileExcerpts, Digest: ref.Digest, IsEgress: true, ManifestOptIn: false})
	if err != ErrFlowDenied {
		t.Fatalf("CheckFlow data class mismatch error = %v, want %v", err, ErrFlowDenied)
	}

	err = store.CheckFlow(FlowCheckRequest{ProducerRole: "workspace", ConsumerRole: "unknown", DataClass: DataClassSpecText, Digest: ref.Digest})
	if err != ErrFlowDenied {
		t.Fatalf("CheckFlow unknown lane error = %v, want %v", err, ErrFlowDenied)
	}

	err = store.CheckFlow(FlowCheckRequest{ProducerRole: "workspace", ConsumerRole: "model_gateway", DataClass: DataClassSpecText, Digest: testDigest("f")})
	if err != ErrArtifactNotFound {
		t.Fatalf("CheckFlow unknown digest error = %v, want %v", err, ErrArtifactNotFound)
	}
}

func TestPromotionRequiresApprovalAndMintsNewReference(t *testing.T) {
	store := newTestStore(t)
	unapproved, err := store.Put(PutRequest{
		Payload:               []byte("sensitive excerpt"),
		ContentType:           "text/plain",
		DataClass:             DataClassUnapprovedFileExcerpts,
		ProvenanceReceiptHash: testDigest("5"),
		CreatedByRole:         "workspace",
	})
	if err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	_, err = store.PromoteApprovedExcerpt(PromotionRequest{UnapprovedDigest: unapproved.Digest})
	if err != ErrPromotionRequiresApproval {
		t.Fatalf("Promote no approver error = %v, want %v", err, ErrPromotionRequiresApproval)
	}
	approved, err := store.PromoteApprovedExcerpt(PromotionRequest{
		UnapprovedDigest:     unapproved.Digest,
		Approver:             "human-1",
		RepoPath:             "repo/file.txt",
		Commit:               "abc123",
		ExtractorToolVersion: "v1",
		FullContentVisible:   true,
	})
	if err != nil {
		t.Fatalf("Promote returned error: %v", err)
	}
	if approved.Digest == unapproved.Digest {
		t.Fatalf("approved digest must differ from unapproved digest")
	}
	if approved.DataClass != DataClassApprovedFileExcerpts {
		t.Fatalf("approved data class = %q", approved.DataClass)
	}
	oldRecord, _ := store.Head(unapproved.Digest)
	if oldRecord.Reference.DataClass != DataClassUnapprovedFileExcerpts {
		t.Fatalf("source artifact mutated: %q", oldRecord.Reference.DataClass)
	}
}

func TestPromotionRateLimitAndBulkGate(t *testing.T) {
	store := newTestStore(t)
	policy := store.Policy()
	policy.MaxPromotionRequestsPerMinute = 1
	if err := store.SetPolicy(policy); err != nil {
		t.Fatalf("SetPolicy error: %v", err)
	}

	first, err := store.Put(PutRequest{Payload: []byte("1"), ContentType: "text/plain", DataClass: DataClassUnapprovedFileExcerpts, ProvenanceReceiptHash: testDigest("6"), CreatedByRole: "workspace"})
	if err != nil {
		t.Fatalf("Put first error: %v", err)
	}
	_, err = store.PromoteApprovedExcerpt(PromotionRequest{UnapprovedDigest: first.Digest, Approver: "human", RepoPath: "a", Commit: "b", ExtractorToolVersion: "c", FullContentVisible: true, BulkRequest: true})
	if err != ErrApprovalBulkConfirmationNeeded {
		t.Fatalf("bulk promotion error = %v, want %v", err, ErrApprovalBulkConfirmationNeeded)
	}
	_, err = store.PromoteApprovedExcerpt(PromotionRequest{UnapprovedDigest: first.Digest, Approver: "human", RepoPath: "a", Commit: "b", ExtractorToolVersion: "c", FullContentVisible: true, BulkRequest: true, BulkApprovalConfirmed: true})
	if err != nil {
		t.Fatalf("bulk promotion confirmed error: %v", err)
	}
	second, err := store.Put(PutRequest{Payload: []byte("2"), ContentType: "text/plain", DataClass: DataClassUnapprovedFileExcerpts, ProvenanceReceiptHash: testDigest("7"), CreatedByRole: "workspace"})
	if err != nil {
		t.Fatalf("Put second error: %v", err)
	}
	_, err = store.PromoteApprovedExcerpt(PromotionRequest{UnapprovedDigest: second.Digest, Approver: "human", RepoPath: "a", Commit: "b", ExtractorToolVersion: "c", FullContentVisible: true})
	if err != ErrPromotionRateLimited {
		t.Fatalf("second promotion error = %v, want %v", err, ErrPromotionRateLimited)
	}
}

func TestQuotasEnforcedAndAudited(t *testing.T) {
	store := newTestStore(t)
	policy := store.Policy()
	policy.PerRoleQuota["workspace"] = Quota{MaxArtifactCount: 1, MaxTotalBytes: 5, MaxSingleArtifactSize: 5}
	if err := store.SetPolicy(policy); err != nil {
		t.Fatalf("SetPolicy error: %v", err)
	}
	_, err := store.Put(PutRequest{Payload: []byte("12345"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("8"), CreatedByRole: "workspace"})
	if err != nil {
		t.Fatalf("first Put error: %v", err)
	}
	_, err = store.Put(PutRequest{Payload: []byte("x"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("8"), CreatedByRole: "workspace"})
	if err != ErrQuotaExceeded {
		t.Fatalf("second Put error = %v, want %v", err, ErrQuotaExceeded)
	}
	audit, err := store.ReadAuditEvents()
	if err != nil {
		t.Fatalf("ReadAuditEvents error: %v", err)
	}
	if !containsAuditType(audit, "artifact_quota_violation") {
		t.Fatalf("expected artifact_quota_violation in audit")
	}
}

func TestRetentionGCAndBackupRestore(t *testing.T) {
	store, keep, backupPath := setupRetentionAndBackupFixture(t)
	assertRetentionAndRestore(t, store, keep, backupPath)
}

func setupRetentionAndBackupFixture(t *testing.T) (*Store, ArtifactReference, string) {
	store, now := setupRetentionStore(t)
	keep := seedRetentionArtifacts(t, store)
	runAndAssertGC(t, store, now, keep)
	backupPath := filepath.Join(t.TempDir(), "backup.json")
	if err := store.ExportBackup(backupPath); err != nil {
		t.Fatalf("ExportBackup error: %v", err)
	}
	return store, keep, backupPath
}

func setupRetentionStore(t *testing.T) (*Store, time.Time) {
	store := newTestStore(t)
	policy := store.Policy()
	policy.UnreferencedTTLSeconds = 1
	if err := store.SetPolicy(policy); err != nil {
		t.Fatalf("SetPolicy error: %v", err)
	}
	now := time.Now().UTC()
	store.nowFn = func() time.Time { return now }
	return store, now
}

func seedRetentionArtifacts(t *testing.T, store *Store) ArtifactReference {
	keep, err := store.Put(PutRequest{Payload: []byte("keep"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("9"), CreatedByRole: "workspace", RunID: "run-active"})
	if err != nil {
		t.Fatalf("Put keep error: %v", err)
	}
	if err := store.SetRunStatus("run-active", "active"); err != nil {
		t.Fatalf("SetRunStatus active error: %v", err)
	}
	if _, err := store.Put(PutRequest{Payload: []byte("drop"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("a"), CreatedByRole: "workspace", RunID: "run-closed"}); err != nil {
		t.Fatalf("Put drop error: %v", err)
	}
	if err := store.SetRunStatus("run-closed", "closed"); err != nil {
		t.Fatalf("SetRunStatus closed error: %v", err)
	}
	return keep
}

func runAndAssertGC(t *testing.T, store *Store, now time.Time, keep ArtifactReference) {
	store.nowFn = func() time.Time { return now.Add(5 * time.Second) }
	gcResult, err := store.GarbageCollect()
	if err != nil {
		t.Fatalf("GarbageCollect error: %v", err)
	}
	if gcResult.FreedBytes == 0 || len(gcResult.DeletedDigests) == 0 {
		t.Fatalf("expected GC to delete at least one artifact")
	}
	if _, err := store.Head(keep.Digest); err != nil {
		t.Fatalf("active run artifact should be retained: %v", err)
	}
}

func assertRetentionAndRestore(t *testing.T, sourceStore *Store, keep ArtifactReference, backupPath string) {
	b, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup error: %v", err)
	}
	var manifest BackupManifest
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatalf("backup json parse error: %v", err)
	}
	if manifest.Schema != "runecode.backup.artifacts.v1" {
		t.Fatalf("backup schema = %q", manifest.Schema)
	}

	restoreStore := newTestStore(t)
	copyBlobsToStore(t, restoreStore, manifest.Artifacts)
	if err := restoreStore.RestoreBackup(backupPath); err != nil {
		t.Fatalf("RestoreBackup error: %v", err)
	}
	if _, err := restoreStore.Head(keep.Digest); err != nil {
		t.Fatalf("restored store missing retained artifact: %v", err)
	}
	_ = sourceStore
}

func TestRestoreRejectsForgedBackupRecord(t *testing.T) {
	store := newTestStore(t)
	ref, err := store.Put(PutRequest{Payload: []byte("payload"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("c"), CreatedByRole: "workspace"})
	if err != nil {
		t.Fatalf("Put error: %v", err)
	}
	backupPath := filepath.Join(t.TempDir(), "backup.json")
	if err := store.ExportBackup(backupPath); err != nil {
		t.Fatalf("ExportBackup error: %v", err)
	}
	b, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("read backup error: %v", err)
	}
	manifest := BackupManifest{}
	if err := json.Unmarshal(b, &manifest); err != nil {
		t.Fatalf("parse backup error: %v", err)
	}
	manifest.Artifacts[0].Reference.Digest = testDigest("d")
	b, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("marshal forged backup error: %v", err)
	}
	if err := os.WriteFile(backupPath, b, 0o644); err != nil {
		t.Fatalf("write forged backup error: %v", err)
	}
	restoreStore := newTestStore(t)
	copyBlobFile(t, store.storeIO.blobPath(ref.Digest), restoreStore.storeIO.blobPath(ref.Digest))
	err = restoreStore.RestoreBackup(backupPath)
	if err == nil {
		t.Fatal("RestoreBackup expected error for forged digest")
	}
}

func TestRestoreRejectsMissingBackupSignature(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.Put(PutRequest{Payload: []byte("payload"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("1"), CreatedByRole: "workspace"}); err != nil {
		t.Fatalf("Put error: %v", err)
	}
	backupPath := filepath.Join(t.TempDir(), "backup.json")
	if err := store.ExportBackup(backupPath); err != nil {
		t.Fatalf("ExportBackup error: %v", err)
	}
	if err := os.Remove(backupSignaturePath(backupPath)); err != nil {
		t.Fatalf("remove signature error: %v", err)
	}
	restoreStore := newTestStore(t)
	err := restoreStore.RestoreBackup(backupPath)
	if err != ErrBackupSignatureMissing {
		t.Fatalf("RestoreBackup error = %v, want %v", err, ErrBackupSignatureMissing)
	}
}

func TestRestoreRejectsTamperedBackupSignature(t *testing.T) {
	store := newTestStore(t)
	if _, err := store.Put(PutRequest{Payload: []byte("payload"), ContentType: "text/plain", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("1"), CreatedByRole: "workspace"}); err != nil {
		t.Fatalf("Put error: %v", err)
	}
	backupPath := filepath.Join(t.TempDir(), "backup.json")
	if err := store.ExportBackup(backupPath); err != nil {
		t.Fatalf("ExportBackup error: %v", err)
	}
	b, err := os.ReadFile(backupSignaturePath(backupPath))
	if err != nil {
		t.Fatalf("read signature error: %v", err)
	}
	sig := BackupSignature{}
	if err := json.Unmarshal(b, &sig); err != nil {
		t.Fatalf("unmarshal signature error: %v", err)
	}
	sig.HMACSHA256 = strings.Repeat("0", len(sig.HMACSHA256))
	b, err = json.MarshalIndent(sig, "", "  ")
	if err != nil {
		t.Fatalf("marshal tampered signature error: %v", err)
	}
	if err := os.WriteFile(backupSignaturePath(backupPath), b, 0o644); err != nil {
		t.Fatalf("write tampered signature error: %v", err)
	}
	restoreStore := newTestStore(t)
	err = restoreStore.RestoreBackup(backupPath)
	if err != ErrBackupSignatureInvalid {
		t.Fatalf("RestoreBackup error = %v, want %v", err, ErrBackupSignatureInvalid)
	}
}

func TestAuditFailureIsSurfaced(t *testing.T) {
	store := newTestStore(t)
	ref, err := store.Put(PutRequest{Payload: []byte("excerpt"), ContentType: "text/plain", DataClass: DataClassUnapprovedFileExcerpts, ProvenanceReceiptHash: testDigest("e"), CreatedByRole: "workspace"})
	if err != nil {
		t.Fatalf("Put error: %v", err)
	}
	badPath := filepath.Join(t.TempDir(), "audit-dir")
	if err := os.MkdirAll(badPath, 0o755); err != nil {
		t.Fatalf("mkdir audit dir error: %v", err)
	}
	store.storeIO.auditPath = badPath
	err = store.CheckFlow(FlowCheckRequest{ProducerRole: "workspace", ConsumerRole: "model_gateway", DataClass: DataClassUnapprovedFileExcerpts, Digest: ref.Digest, IsEgress: true})
	if err == nil {
		t.Fatal("CheckFlow expected audit write error")
	}
}

func TestCanonicalJSONRejectsNonIntegerNumbers(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Put(PutRequest{Payload: []byte(`{"a":1.2}`), ContentType: "application/json", DataClass: DataClassSpecText, ProvenanceReceiptHash: testDigest("1"), CreatedByRole: "workspace"})
	if err == nil {
		t.Fatal("Put expected canonicalization error for non-integer JSON")
	}
}

func TestReservedDataClassesFailClosedByDefault(t *testing.T) {
	store := newTestStore(t)
	_, err := store.Put(PutRequest{Payload: []byte("web"), ContentType: "text/plain", DataClass: DataClassWebQuery, ProvenanceReceiptHash: testDigest("b"), CreatedByRole: "workspace"})
	if err != ErrReservedDataClassDisabled {
		t.Fatalf("reserved class Put error = %v, want %v", err, ErrReservedDataClassDisabled)
	}
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	t.Setenv(backupHMACKeyEnv, "test-backup-key")
	root := t.TempDir()
	store, err := NewStore(root)
	if err != nil {
		t.Fatalf("NewStore returned error: %v", err)
	}
	store.nowFn = func() time.Time { return time.Now().UTC() }
	return store
}

func testDigest(seed string) string {
	base := strings.Repeat(seed, 64)
	if len(base) > 64 {
		base = base[:64]
	}
	for len(base) < 64 {
		base += "0"
	}
	return "sha256:" + base
}

func containsAuditType(events []AuditEvent, eventType string) bool {
	for _, event := range events {
		if event.Type == eventType {
			return true
		}
	}
	return false
}

func ioReadAllAndClose(r io.ReadCloser) ([]byte, error) {
	b, err := io.ReadAll(r)
	_ = r.Close()
	return b, err
}

func copyBlobsToStore(t *testing.T, dst *Store, records []ArtifactRecord) {
	t.Helper()
	for _, rec := range records {
		copyBlobFile(t, rec.BlobPath, dst.storeIO.blobPath(rec.Reference.Digest))
	}
}

func copyBlobFile(t *testing.T, src, dst string) {
	t.Helper()
	b, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read blob %s error: %v", src, err)
	}
	if err := os.WriteFile(dst, b, 0o644); err != nil {
		t.Fatalf("write blob %s error: %v", dst, err)
	}
}
