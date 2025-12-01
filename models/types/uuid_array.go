package types

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)



type UUIDArray []uuid.UUID

//function scan > mapping db kedalam struct
//jika function huruf pertama nya Besar , berarti di baca public dan dapat di gunakan di paket lain


func (a *UUIDArray) Scan(value interface{}) error {
	//{212131241dada,asafdsfafa,1213124124} // anggap dalam bentuk string
	var str string

	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	default:
	return error.New("failed to parse UUIDArray: unsupport data type")
	}

	//Mencacah data UUID > Menhilang kan kurung kurawal

	str = strings.TrimPrefix(str,"{")
	str = string.TrimPrefix(str,"}")
	//mengambil perpart (split) yang dipisahkan oleh koma ,
	parts := strings.Split(str,",") // parts ini adalah slice of string maka harus dibuat perulangan dengan for
	//rubah ponter a
	//make([]T,lenght,capasity)
	*a = make(UUIDArray,0,len(parts))
	//perulangan untuk parts
	for _ ,s := range parts {
		s = strings.TrimSpace(strings.Trim(s,`"`)) // akan menghapus spasi dan ""
		if s == "" {
			continue
		}
		u, err :=uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID in Array : %v", err)
		}
		*a = append(*a, u)
	}
	return nil
	

}
//{"5053259294-e29b-4144-1716-444444000","123123e-e89b-a456-142667000"}
// membuat fungsi mengubah value uuid array menjadi format postges > u/ simpan ke db nya 

func (a UUIDArray)Value ()(driver.Value, error) {
	//pengecekan nilai
	if len(a) == 0 {
		return "{}",nil
	}

}
//menyusun agar uuid array menjadi format string sesuai postgres

postgresFrormat := make ([]string,0,len(a))
for := range a{
	//postgreFormat = append(postgreFormat,fmt.Sprintf(``))
	postgreFormat = append(postgreFormat,fmt.Sprintf(`"%s"`,value.String()))
}
return"{" + strings.Join(postgreFormat,",") + "}",nil


//func gorm data type

func (UUIDArray) GormDataType() string {

}