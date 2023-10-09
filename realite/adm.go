package main

// func proxyConnect(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("CONNECT requested to %v (from %v)", req.Host, req.RemoteAddr)
// 	targetConn, err := net.Dial("tcp", req.Host)
// 	if err != nil {
// 		log.Println("failed to dial to target", req.Host)
// 		http.Error(w, err.Error(), http.StatusServiceUnavailable)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	hj, ok := w.(http.Hijacker)
// 	if !ok {
// 		fmt.Println("http server doesn't support hijacking connection")
// 		return
// 	}

// 	clientConn, _, err := hj.Hijack()
// 	if err != nil {
// 		fmt.Println("http hijacking failed")
// 		return
// 	}

// 	log.Println("tunnel established")
// 	var g errgroup.Group
// 	g.Go(func() error {
// 		defer targetConn.Close()
// 		_, err := io.Copy(targetConn, clientConn)
// 		return err
// 	})
// 	g.Go(func() error {
// 		defer clientConn.Close()
// 		_, err := io.Copy(clientConn, targetConn)
// 		return err
// 	})
// 	g.Wait()
// }
