package blob

type RetrieveBlob struct{

}

func NewRetrieveBlob() *RetrieveBlob{
	return &RetrieveBlob{}
}

func (r *RetrieveBlob) retrieveBlob(path string) string{
	return ""
}