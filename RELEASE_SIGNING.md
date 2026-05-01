# Release Signing

This project uses GitHub Actions to build Tauri release artifacts. macOS releases require Apple Developer ID signing and notarization so the downloaded app can pass Gatekeeper checks.

The release workflow reads the following GitHub Actions secrets during macOS builds:

```text
APPLE_CERTIFICATE
APPLE_CERTIFICATE_PASSWORD
APPLE_SIGNING_IDENTITY
APPLE_ID
APPLE_PASSWORD
APPLE_TEAM_ID
```

Configure them in GitHub:

```text
Repository -> Settings -> Secrets and variables -> Actions -> New repository secret
```

Do not commit certificate files, passwords, or generated secret values to the repository.

## Required Apple Resources

You need an active Apple Developer Program membership.

Create or use a certificate of this type:

```text
Developer ID Application
```

This is the certificate type used for distributing macOS apps outside the Mac App Store. Do not use `Apple Development`, `Mac App Distribution`, or `Developer ID Installer` for the Tauri app bundle signing step.

The app bundle identifier is currently configured in `src-tauri/tauri.conf.json`:

```json
"identifier": "com.is14w.sekaitext"
```

Make sure this identifier is available under the Apple Developer team used for signing.

## Secrets

### `APPLE_CERTIFICATE`

Base64-encoded `.p12` export of the `Developer ID Application` certificate and its private key.

Export it from macOS Keychain Access:

1. Open Keychain Access.
2. Find the `Developer ID Application: ...` certificate.
3. Expand the certificate and confirm it includes the private key.
4. Select the certificate and private key.
5. Export as `.p12`.
6. Set an export password.

Then encode the `.p12` file:

```bash
base64 -i DeveloperIDApplication.p12 | pbcopy
```

Paste the copied base64 value into the GitHub secret.

### `APPLE_CERTIFICATE_PASSWORD`

The password used when exporting the `.p12` certificate from Keychain Access.

Example value format:

```text
your-p12-export-password
```

### `APPLE_SIGNING_IDENTITY`

The exact signing identity name for the Developer ID Application certificate.

Example value format:

```text
Developer ID Application: Your Name (ABCDE12345)
```

You can find the identity locally with:

```bash
security find-identity -v -p codesigning
```

Use the `Developer ID Application` identity that matches the Apple Developer team.

### `APPLE_ID`

The Apple ID email address used for notarization.

Example value format:

```text
name@example.com
```

This Apple ID must have access to the Apple Developer team used for signing.

### `APPLE_PASSWORD`

An app-specific password for the Apple ID used by notarization.

Do not use the normal Apple ID login password.

Create it from the Apple ID account page:

```text
https://account.apple.com/account/manage
```

Use the generated app-specific password as the secret value.

Example value format:

```text
abcd-efgh-ijkl-mnop
```

### `APPLE_TEAM_ID`

The Apple Developer Team ID.

Example value format:

```text
ABCDE12345
```

You can find it in the Apple Developer account membership details or in the certificate identity string.

## Verification

After all secrets are configured, create a tag release or manually run the release workflow.

The macOS jobs should run the `Build and sign Tauri app` step. If signing or notarization fails, check the GitHub Actions logs for the Apple certificate import, signing identity, or notarization error.

Common issues:

- `APPLE_CERTIFICATE` is not valid base64.
- The `.p12` export does not include the private key.
- `APPLE_CERTIFICATE_PASSWORD` does not match the `.p12` export password.
- `APPLE_SIGNING_IDENTITY` does not exactly match the certificate identity.
- `APPLE_PASSWORD` is the Apple ID login password instead of an app-specific password.
- `APPLE_TEAM_ID` does not match the certificate team.
