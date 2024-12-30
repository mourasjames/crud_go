package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Equipamento struct {
	ID   string
	NOME string
	IP   string
	RACK int
}

func NovoEquipamento(nome string, ip string, rack int) *Equipamento {
	return &Equipamento{
		ID:   uuid.New().String(),
		NOME: nome,
		IP:   ip,
		RACK: rack,
	}
}

func main() {

	db, err := sql.Open("mysql", "go:goexpert@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	equipamento, err := SelecaoDb(db)
	if err != nil {
		panic(err)
	}
	fmt.Println("EQUIPAMENTOS CADASTRADOS: ")
	fmt.Println("----------")
	fmt.Println("ID  -  Nome  -  IP  -  RACK  ")
	for _, e := range equipamento {
		fmt.Println(e.ID, " - ", e.NOME, " - ", e.IP, " - ", e.RACK)
	}

	fmt.Println()
	fmt.Println("Quais das opções abaixo:")
	fmt.Println(" 1 - Lista item no DB;")
	fmt.Println(" 2 - Inserir dados no DB;")
	fmt.Println(" 3 - Alterar dados no DB;")
	fmt.Println(" 4 - Remover dados no DB;")
	fmt.Println(" 5 - Sair.")

	var opcao int
	_, err = fmt.Scan(&opcao)
	if err != nil {
		panic(err)
	}

	switch opcao {
	case 1:
		fmt.Printf("Qual o ID do item: ")
		var itemid string
		fmt.Scan(&itemid)
		if err != nil {
			panic(err)
		}
		e, err := selecaoItem(db, itemid)
		if err != nil {
			panic(err)
		}
		fmt.Println("----")
		fmt.Println("Item selecionado: " /**e.ID, "-", **/, e.NOME, "-", e.IP, "-", e.RACK)
		fmt.Println("----")
	case 2:
		var e Equipamento
		fmt.Println("Digite o ID do equipamento: ")
		fmt.Scan(&e.ID)
		fmt.Println("Digite o NOME do equipamento: ")
		fmt.Scan(&e.NOME)
		fmt.Println("Digite o IP do equipamento: ")
		fmt.Scan(&e.IP)
		fmt.Println("Digite o RACK do equipamento: ")
		fmt.Scan(&e.RACK)

		err := inserirDados(db, e)
		if err != nil {
			log.Printf("Erro ao iserir item: %v\n", err)
		}
		fmt.Println("Equimento inserido com sucesso!")
	case 3:
		fmt.Println("Alterar qual ID: ")
		var itemID string
		_, err := fmt.Scan(&itemID)
		if err != nil {
			panic(err)
		}
		a, err := selecaoItem(db, itemID)
		if err != nil {
			fmt.Printf("Item inesistente: %v\n", err)
			break
		}
		fmt.Println("Digite o NOME do equipamento: ")
		fmt.Scan(&a.NOME)
		fmt.Println("Digite o IP do equipamento: ")
		fmt.Scan(&a.IP)
		fmt.Println("Digite o RACK do equipamento: ")
		fmt.Scan(&a.RACK)

		err = alterarDados(db, *a)
		if err != nil {
			log.Printf("Erro na alteração: %v\n", err)
		}
		fmt.Println("Equimento Alterado com sucesso!")
	case 4:
		fmt.Println("Remover qual ID: ")
		var itemID string
		_, err := fmt.Scan(&itemID)
		if err != nil {
			panic(err)
		}
		a, err := selecaoItem(db, itemID)
		if err != nil {
			fmt.Printf("Item inesistente: %v\n", err)
			break
		}
		var esc string
		fmt.Println("Remover ID: (s) ou (n)")
		fmt.Scan(&esc)
		if esc == "s" {
			err := removerDados(db, a.ID)
			if err != nil {
				log.Printf("Erro no remover: %v\n", err)
				break
			}
			fmt.Println("Equimento removido com sucesso!")
		} else {
			fmt.Println("Exclusão cancelada!")
		}

	case 5:
		fmt.Println("Bye...")
		break
	default:
		fmt.Println("Opcão inválida!")
		main()
	}
}

func SelecaoDb(db *sql.DB) ([]Equipamento, error) {
	sel, err := db.Query("select id, nome, ip, rack from equipamento")
	if err != nil {
		return nil, err
	}
	defer sel.Close()

	var equipamento []Equipamento
	for sel.Next() {
		var e Equipamento
		err := sel.Scan(&e.ID, &e.NOME, &e.IP, &e.RACK)
		if err != nil {
			return nil, err
		}
		equipamento = append(equipamento, e)
	}
	return equipamento, nil
}

func selecaoItem(db *sql.DB, itemid string) (*Equipamento, error) {
	sel, err := db.Prepare("select id, nome, ip, rack from equipamento where id = ?")
	if err != nil {
		return nil, err
	}
	defer sel.Close()
	var e Equipamento
	err = sel.QueryRow(itemid).Scan(&e.ID, &e.NOME, &e.IP, &e.RACK)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func inserirDados(db *sql.DB, equip Equipamento) error {
	inserir, err := db.Prepare("insert into equipamento(id, nome, ip, rack) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer inserir.Close()
	_, err = inserir.Exec(equip.ID, equip.NOME, equip.IP, equip.RACK)
	if err != nil {
		return err
	}
	return nil
}

func alterarDados(db *sql.DB, alter Equipamento) error {
	alt, err := db.Prepare("update equipamento set nome = ?, ip = ?, rack = ? where id = ?")
	if err != nil {
		return err
	}
	defer alt.Close()
	_, err = alt.Exec(alter.NOME, alter.IP, alter.RACK, alter.ID)
	if err != nil {
		return err
	}
	return nil
}

func removerDados(db *sql.DB, itemID string) error {
	rem, err := db.Prepare("delete from equipamento where id = ?")
	if err != nil {
		return err
	}
	_, err = rem.Exec(itemID)
	if err != nil {
		return err
	}
	return nil
}
