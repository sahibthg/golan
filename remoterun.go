func remoteRun(host,cmd,secret string ) string {
        user := os.Getenv("LOGNAME")
        whoami := os.Getenv("USER")
        command := cmd+" "+secret

        key, err := ioutil.ReadFile("/home/"+whoami+"/.ssh/id_rsa")
        if err != nil {
                log.Fatalf("unable to read private key: %v", err)
        }

        // Create the Signer for this private key.
        signer, err := ssh.ParsePrivateKey(key)
        if err != nil {
                log.Fatalf("unable to parse private key: %v", err)
        }

        hostKeyCallback, err := kh.New("/home/"+whoami+"/.ssh/known_hosts")
        if err != nil {
                log.Fatal("could not create hostkeycallback function: ", err)
        }

        config := &ssh.ClientConfig{
                User: user,
                Auth: []ssh.AuthMethod{
                        // Add in password check here for moar security.
                        ssh.PublicKeys(signer),
                },
                HostKeyCallback: hostKeyCallback,
        }
        // Connect to the remote server and perform the SSH handshake.
        client, err := ssh.Dial("tcp", host+":22", config)
        if err != nil {
                log.Fatalf("unable to connect: %v", err)
        }
        defer client.Close()
        ss, err := client.NewSession()
        if err != nil {
                log.Fatal("unable to create SSH session: ", err)
        }
        defer ss.Close()
        // Creating the buffer which will hold the remotly executed command's output.
        var stdoutBuf bytes.Buffer
        ss.Stdout = &stdoutBuf
        ss.Run(command)
        // Let's print out the result of command.
        fmt.Println(stdoutBuf.String())
        return stdoutBuf.String()
}
