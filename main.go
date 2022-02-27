/*
Copyright © 2022 Quantos developers

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"Quantos/sdk"
	"log"
)

func main() {
	/*c, err := config.Init(os.Stdout)
	addrSDK := sdk.GetAddressSDK()
	addrSDK.InitSDK("test")
	var out string
	o := addrSDK.GenerateQBITWalletAddress(out)
	z := addrSDK.GetZeroAddress()
	spew.Dump(z)
	fmt.Printf("New Wallet Address (QBIT): %s \n", o)
	if err != nil {
		panic(err)
	}
	spew.Dump(c.Identity)*/

	//cmd.Execute()
	raw, m := sdk.GetAddressSDK().GenerateMasterWalletAddress()
	d := sdk.GetAddressSDK().DeriveFromMaster(raw, m)

	log.Printf("Master Key (to unlock your wallet): %s", m)
	log.Printf("Your wallet address (long form): %s \n", d[:40])

}
