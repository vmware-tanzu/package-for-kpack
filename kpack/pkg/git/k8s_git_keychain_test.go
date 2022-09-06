package git

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	git2go "github.com/libgit2/git2go/v33"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	buildapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha2"
)

func Test(t *testing.T) {
	privateKeyBytes := gitTest{key1: generateRandomPrivateKey(t), key2: generateRandomPrivateKey(t)}
	spec.Run(t, "Test Git Keychain", privateKeyBytes.testK8sGitKeychain)
}

type gitTest struct {
	key1 []byte
	key2 []byte
}

func (keys gitTest) testK8sGitKeychain(t *testing.T, when spec.G, it spec.S) {
	const serviceAccount = "some-service-account"
	const testNamespace = "test-namespace"

	var (
		fakeClient = fake.NewSimpleClientset(
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-1",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://github.com",
					},
				},
				Type: v1.SecretTypeBasicAuth,
				Data: map[string][]byte{
					v1.BasicAuthUsernameKey: []byte("saved-username"),
					v1.BasicAuthPasswordKey: []byte("saved-password"),
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-2",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "noschemegit.com",
					},
				},
				Type: v1.SecretTypeBasicAuth,
				Data: map[string][]byte{
					v1.BasicAuthUsernameKey: []byte("noschemegit-username"),
					v1.BasicAuthPasswordKey: []byte("noschemegit-password"),
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-3",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://bitbucket.com",
					},
				},
				Type: v1.SecretTypeSSHAuth,
				Data: map[string][]byte{
					v1.SSHAuthPrivateKey: keys.key1,
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-4",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://gitlab.com",
					},
				},
				Type: v1.SecretTypeSSHAuth,
				Data: map[string][]byte{
					v1.SSHAuthPrivateKey: keys.key1,
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-5",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://gitlab.com",
					},
				},
				Type: v1.SecretTypeBasicAuth,
				Data: map[string][]byte{
					v1.BasicAuthUsernameKey: []byte("gitlab-username"),
					v1.BasicAuthPasswordKey: []byte("gitlab-password"),
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-6",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://gitlab.com",
					},
				},
				Type: v1.SecretTypeSSHAuth,
				Data: map[string][]byte{
					v1.SSHAuthPrivateKey: keys.key2,
				},
			},
			&v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret-7",
					Namespace: testNamespace,
					Annotations: map[string]string{
						buildapi.GITSecretAnnotationPrefix: "https://github.com",
					},
				},
				Type: v1.SecretTypeBasicAuth,
				Data: map[string][]byte{
					v1.BasicAuthUsernameKey: []byte("other-username"),
					v1.BasicAuthPasswordKey: []byte("other-password"),
				},
			},
			&v1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      serviceAccount,
					Namespace: testNamespace,
				},
				Secrets: []v1.ObjectReference{
					{Name: "secret-1"},
					{Name: "secret-2"},
					{Name: "secret-3"},
					{Name: "secret-4"},
					{Name: "secret-5"},
					{Name: "secret-6"},
					{Name: "secret-7"},
				},
			})
		keychainFactory = newK8sGitKeychainFactory(fakeClient)
	)

	when("Keychain resolves", func() {
		var keychain GitKeychain

		it.Before(func() {
			var err error
			keychain, err = keychainFactory.KeychainForServiceAccount(context.Background(), testNamespace, serviceAccount)
			require.NoError(t, err)
		})

		it("returns  alphabetical first git Auth for matching secrets with basic auth", func() {
			cred, err := keychain.Resolve("https://github.com/org/repo", "", git2go.CredentialTypeUserpassPlaintext)
			require.NoError(t, err)

			require.Equal(t, BasicGit2GoAuth{
				Username: "saved-username",
				Password: "saved-password",
			}, cred)

			git2goCred, err := cred.Cred()
			require.NoError(t, err)

			require.Equal(t, git2goCred.Type(), git2go.CredentialTypeUserpassPlaintext)
		})

		it("returns the alphabetical first secretRef for ssh auth", func() {
			cred, err := keychain.Resolve("https://gitlab.com/my-repo.git", "gituser", git2go.CredentialTypeSSHKey)
			require.NoError(t, err)

			require.Equal(t, SSHGit2GoAuth{
				Username:   "gituser",
				PrivateKey: string(keys.key1),
			}, cred)

			git2goCred, err := cred.Cred()
			require.NoError(t, err)

			require.Equal(t, git2goCred.Type(), git2go.CredentialTypeSSHCustom)
		})

		it("returns git Auth for matching secrets without scheme", func() {
			cred, err := keychain.Resolve("https://noschemegit.com/org/repo", "", git2go.CredentialTypeUserpassPlaintext)
			require.NoError(t, err)

			require.Equal(t, BasicGit2GoAuth{
				Username: "noschemegit-username",
				Password: "noschemegit-password",
			}, cred)
		})

		it("returns an error if no credentials found", func() {
			_, err := keychain.Resolve("https://no-creds-github.com/org/repo", "git", git2go.CredentialTypeUserpassPlaintext)
			require.EqualError(t, err, "no credentials found for https://no-creds-github.com/org/repo")
		})
	})
}

func generateRandomPrivateKey(t *testing.T) []byte {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)
	var pemBlock = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(pemBlock)
}
