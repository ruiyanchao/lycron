package lycron

type Store interface {
	Save()
	GetRow()
	GetRows()
	GetList()
	Delete()
}

type FileStore struct {

}

func(fs *FileStore)Save(){

}

func(fs *FileStore)GetRow(){

}

func(fs *FileStore)GetRows(){

}

func(fs *FileStore)GetList(){

}

func(fs *FileStore)Delete(){

}

