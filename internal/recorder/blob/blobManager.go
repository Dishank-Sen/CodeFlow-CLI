package blob

type ISaveBlob interface{
	saveBlobFromPath(hash string, path string) string
	saveBlobFromContent(hash string, content string) string
}

type IRetrieveBlob interface{
	retrieveBlob(path string) string
}
