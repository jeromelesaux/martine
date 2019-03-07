# martine

After Claudia from Eliot a long time ago, here is coming Martine.
Convert JPEG or PNG image file into  SCR/PAL or Overscan  Amstrad CPC file screen.
Multi os, you can convert any pictures to Amstrad CPC Screen.

To Install and compile
```go get github.com/jeromelesaux/martine
cd $GOPATH/src/github.com/jeromelesaux/martine
go get 
go build
```
Usage and options : 

```martine to convert image to Amstrad cpc screen (even overscan)
By Impact Sid (Version:0.1Beta)
usage :

  -a int
    	Algorithm to resize the image (available 1: NearestNeighbor (default), 2: CatmullRom, 3: Lanczos, 4: Linear) (default 1)
  -f	Overscan mode (default no overscan)
  -h int
    	Custom output height in pixels. (default -1)
  -m string
    	Output mode to use (mode0,mode1,mode2 or overscan available).
  -o string
    	Output directory
  -p string
    	Picture path of the Amsdos file.
  -w int
    	Custom output width in pixels. (default -1)
```

examples :

* convert samples/Batman-Neal-Adams.jpg 

  * in mode 0 
```martine -p samples/Batman-Neal-Adams.jpg -m 0```
  * in mode 1 
```martine -p samples/Batman-Neal-Adams.jpg -m 1```
  * in mode 2 
```martine -p samples/Batman-Neal-Adams.jpg -m 2```
  * in mode 0 in overscan : 
```martine -p samples/Batman-Neal-Adams.jpg -m 0 -f```

