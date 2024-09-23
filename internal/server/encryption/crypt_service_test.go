package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CryptServiceTestSuite struct {
	suite.Suite
	cryptService *Service
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CryptServiceTestSuite))
}

func (c *CryptServiceTestSuite) SetupSuite() {
	service, err := New([]byte("test-master-key"))
	require.NoError(c.T(), err)
	c.cryptService = service
}

func (c *CryptServiceTestSuite) Test_EncryptDecryptWithMasterKey() {
	data := []byte("hello world")
	cryptData, err := c.cryptService.EncryptWithMasterKey(data)
	assert.NoError(c.T(), err)
	assert.NotEmpty(c.T(), cryptData)
	decryptedData, err := c.cryptService.DecryptWithMasterKey(cryptData)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), data, decryptedData)
}

func (c *CryptServiceTestSuite) Test_GenerateKey() {
	key, err := c.cryptService.GenerateKey()
	assert.NoError(c.T(), err)
	assert.NotEmpty(c.T(), key)
	assert.Equal(c.T(), 32, len(key))
}

func (c *CryptServiceTestSuite) Test_EncryptDecrypt() {
	key, err := c.cryptService.GenerateKey()
	assert.NoError(c.T(), err)
	assert.NotEmpty(c.T(), key)
	assert.Equal(c.T(), 32, len(key))

	data := []byte("hello world")
	cryptData, err := c.cryptService.Encrypt(key, data)
	assert.NoError(c.T(), err)
	assert.NotEmpty(c.T(), cryptData)
	decryptedData, err := c.cryptService.Decrypt(key, cryptData)
	assert.NoError(c.T(), err)
	assert.Equal(c.T(), data, decryptedData)
}
