package mydict

import "errors"

type Dictionary map[string] string

var (
errNotFound=errors.New("Not Found!")
errCantUpdate=errors.New("Cant update non-existing word")
errCantDelete=errors.New("Cant delete non-existing word")
errWordExists=errors.New("That word already exists.")
)

//Search a word at the dictionary
func (d Dictionary) Search(word string) (string,error){
	value,exist :=d[word]
	if exist {
		return value,nil
	}
	return "",errNotFound
}

//Add a word to the dictionary
func (d Dictionary) Add(word, def string) (error) {
	_,err := d.Search(word)
	switch err {
	case errNotFound:
		d[word]=def
	case nil:
		return errWordExists
	}
	return nil
}

//Update a word
func (d Dictionary) Update(word,def string) (error){
	_,err:=d.Search(word)
	switch err {
	case errNotFound:
		return errCantUpdate
	case nil:
		d[word]=def
	}
	return nil
}

//Delete a word
func (d Dictionary) Delete(word string) (error){
	_,err:=d.Search(word)
	switch err {
	case errNotFound:
		return errCantDelete
	case nil:
		delete(d,word)
	}
	return nil
}