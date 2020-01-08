package main

import (
    "fmt"
    "os"
    "flag"
    "os/user"
   // "os/exec"
    //"io/ioutil"
    "bufio"
    "strings"
    "github.com/olekukonko/tablewriter"
    "errors"
)

type Alias struct {
    name string
    command string
}

type AliasSsh struct {
    alias string
    username string
    host string
    password string
}

func main() {
    arguments := os.Args[1:]

    var flagCreateType string

    list := flag.Bool("l", false, "List aliases")

    flag.StringVar(&flagCreateType, "create", "", "Type entity")
    flag.StringVar(&flagCreateType, "c", "", "Type entity"+" (shorthand)")

    flag.Parse()

    if flagCreateType == "ssh" {
        doCreateSshAlias()
        return
    }

    if len(arguments) == 0 || *list == true {
        displayListAliases()
        return
    }

    if len(arguments) == 1 {
        doCreateSimpleAlias(arguments[0])
    }

    //fmt.Println(arguments)

}

func doCreateSimpleAlias(cmd string) {
    chunks := strings.Split(cmd, "=")
    if len(chunks) < 2 {
        panic("For create alias, must type aliasName=command")
    }
    var command string
    command = fmt.Sprintf("alias %s='%s'", chunks[0], chunks[1])
    appendAliasToFile(command)
    fmt.Println(command)
}

func doCreateSshAlias() {
    var aliasCmd string

    ssh, err := getSshUserData()
    if err != nil {
        panic(err)
    }

    aliasCmd = fmt.Sprintf(
        "alias %s='sshpass -p %s ssh -o StrictHostKeyChecking=no %s@%s'\n",
        trim(ssh.alias),
        trim(ssh.password),
        trim(ssh.username),
        trim(ssh.host))

    appendAliasToFile(aliasCmd)

    fmt.Println(aliasCmd)
}

func appendAliasToFile(command string) {
    var filePath string = getAliasFilePath()
    f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        panic(err)
    }
    defer f.Close()

    if _, err = f.WriteString(command);
    err != nil {
        panic(err)
    }
}

func getSshUserData() (AliasSsh, error) {
    var ssh AliasSsh

    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter Alias Name: ")
    ssh.alias, _ = reader.ReadString('\n')

    if len(ssh.alias) <= 1 {
        return ssh, errors.New("Alias must be not empty")
    }

    fmt.Print("Enter Username: ")
    ssh.username, _ = reader.ReadString('\n')
    if len(ssh.username) <= 1 {
        return ssh, errors.New("Username must be not empty")
    }

    fmt.Print("Enter Host: ")
    ssh.host, _ = reader.ReadString('\n')
    //XXX: add condition IP address
    if len(ssh.host) <= 1 {
        return ssh, errors.New("Host must be not empty")
    }

    fmt.Print("Enter password: ")
    ssh.password, _ = reader.ReadString('\n')
    if len(ssh.password) <= 1 {
        return ssh, errors.New("Password must be not empty")
    }

    return ssh, nil
}

func trim(str string) (string) {
    return strings.Trim(str, " \r\n");
}

func displayListAliases() {
    var alias Alias
    var filePath string;
    filePath = getAliasFilePath()

    file, err := os.Open(filePath)
    check(err)
    defer file.Close()

    scanner := bufio.NewScanner(file)

    data := [][]string{}

    for scanner.Scan() {

        arr := strings.Split(scanner.Text(), " ")

        if arr[0] == "alias" {
            alias = getAliasStruct(arr)

            data  = append(data, []string{alias.name, alias.command})
        }
    }

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Name", "Command"})
    //table.SetRowLine(true)

    table.SetAutoWrapText(false)

    table.SetColumnColor(
        tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor},
    	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor})


    table.AppendBulk(data)
    table.Render()
}

func getAliasStruct(arr []string) (Alias) {
    var aliasString, command string

    aliasString = arrayToString(arr, " ")

    aliasCommand := strings.Split(aliasString, "=")

    command = arrayToString(aliasCommand, "=")

    return Alias{aliasCommand[0], command}
}

func arrayToString(data []string, separator string, params ...int) (string) {
    var key int = 1;

    if len(params) > 0 {
        key = params[0]
    }

    return strings.Join(data[key:], separator)
}

func getHomeDir() (string) {
    usr, _ := user.Current()
    return usr.HomeDir
}

func getAliasFilePath() (string) {
    //var homeDir string = getHomeDir()
    //fmt.Println(homeDir + "/.bash_profile")
    return "/Users/nil/.bash_profile";
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}