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

```martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)
By Impact Sid (Version:0.1.Alpha)
usage :

  -a int
    	Algorithm to resize the image (available :
    		1: NearestNeighbor (default)
    		2: CatmullRom
    		3: Lanczos,
    		4: Linear
    		5: Box
    		6: Hermite
    		7: BSpline
    		8: Hamming
    		9: Hann
    		10: Gaussian
    		11: Blackman
    		12: Bartlett
    		13: Welch
    		14: Cosine (default 1)
  -f	Overscan mode (default no overscan)
  -h int
    	Custom output height in pixels. (default -1)
  -help
    	Display help message
  -i string
    	Picture path of the input file.
  -m int
    	Output mode to use :
    		0 for mode0
    		1 for mode1
    		2 for mode2
    		and add -f option for overscan export. (default -1)
  -n	no amsdos header for all files (default amsdos header added).
  -o string
    	Output directory
  -p	Plus mode (means generate an image for CPC Plus Screen)
  -s string
    	Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb
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

