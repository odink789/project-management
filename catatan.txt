1. https://github.com/odink789/project-management
2. go get github.com/gofiber/fiber/v2           << install framework


user postgres : postgres passwd :111213aa

Model ?

Model adalah representasi data yang kita simpan didb . ibarat nya model adalah
blueprint atau cetak bitu untuk data

analogi :

bayangkan model seperti ktp - ktp mendeskripsikan siapa kamu (nama,ttl ,alamat ,), tapi ktp itu sendiri
bukan kamu kan .
begitu pula model - model hanya mendeskripsikan bentuk data , bukan data itu sendiri

contoh model :

Type User struct {
    ID int 
    Name string
    Email string
    Password string
}

Relasi adalah Cara menghubungkan data dengan model nya 

Relasi ada 3

one to one > satu user -> satu profile > hubungan ekslusif
one to many > satu board > banyak list > hubungan paling umum
many to many > banyak user > banyak board  >>biasa nya perlu pivot table

UUID = Universally Unique IDentifier
 adalah angka unik 128bit(panjangnya 16 byte) yg digunakan untuk mengidentifikasi sesuatu (record,user,file,dsb) di seluruh dunia tanpa tabrakan(collison)
UUID sering digunakan sebagai primary key atau identifier unik, terutama ketika :
* Data dihasilkan dari banyak sistem berbeda
* Tidak ingin bergantung pada auto-increment ID (1,2,3...)
* Perlu ID yang unik secara global

