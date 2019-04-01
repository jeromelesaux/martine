# martine

After Claudia from Eliot a long time ago, here is coming Martine.
A cli which converts JPEG or PNG image file into  SCR/PAL or Overscan  Amstrad CPC file screen.
Multi os, you can convert any pictures to Amstrad CPC Screen.
The files generated (.win, .scr, .ink) are compatible with [OCP art studio](http://www.cpc-power.com/index.php?page=detail&num=4963) and Impdraw V2 [i2](http://amstradplus.forumforever.com/t462-iMPdraw-v2-0.htm)

To Install and compile
```go get github.com/jeromelesaux/martine
cd $GOPATH/src/github.com/jeromelesaux/martine
go get 
go build```

To get binary : 
[https//github.com/jeromelesaux/martine/releases](https//github.com/jeromelesaux/martine/releases)
<br>OS avaible : Linux, Macos X and Windows  

Usage and options : 

```martine convert (jpeg, png format) image to Amstrad cpc screen (even overscan)
By Impact Sid (Version:0.3)
Special thanks to @Ast (for his support), @Siko and @Tronic for ideas
usage :

  -a int
        Algorithm to resize the image (available : 
                1: NearestNeighbor (default)
                2: CatmullRom
                3: Lanczos
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
                14: Cosine
                 (default 1)
  -f    Overscan mode (default no overscan)
  -h int
        Custom output height in pixels. (default -1)
  -help
        Display help message
  -i string
        Picture path of the input file.
  -iter int
        Iterations number to walk in tile mode (default -1)
  -keephigh int
        bit rotation on the top and keep pixels (default -1)
  -keeplow int
        bit rotation on the bottom and keep pixels (default -1)
  -losthigh int
        bit rotation on the top and lost pixels (default -1)
  -lostlow int
        bit rotation on the bottom and lost pixels (default -1)
  -m int
        Output mode to use :
                0 for mode0
                1 for mode1
                2 for mode2
                and add -f option for overscan export.
                 (default -1)
  -n    no amsdos header for all files (default amsdos header added).
  -o string
        Output directory
  -p    Plus mode (means generate an image for CPC Plus Screen)
  -rla int
        bit rotation on the left and keep pixels (default -1)
  -roll
        Roll mode allow to walk and walk into the input file.
  -rra int
        bit rotation on the right and keep pixels (default -1)
  -s string
        Byte statement to replace in ascii export (default is BYTE), you can replace or instance by defb
  -sla int
        bit rotation on the left and lost pixels (default -1)
  -sra int
        bit rotation on the right and lost pixels (default -1)
  -w int
        Custom output width in pixels. (default -1)
```

examples :

* convert samples/Batman-Neal-Adams.jpg 

  * in mode 0 
```martine -i samples/Batman-Neal-Adams.jpg -m 0```
  * in mode 1 
```martine -i samples/Batman-Neal-Adams.jpg -m 1```
  * in mode 2 
```martine -i samples/Batman-Neal-Adams.jpg -m 2```
  * in mode 0 in overscan : 
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f```
  * in mode 0 overscan for Plus series :
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f -p```
  * to get sprites (40 pixels wide)
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -w 40```
  * roll mode to do an rra operation on the image (will create 16 sprites with a rra operation on the first pixels on the left)
	```martine -i samples/Batman-Neal-Adams.jpg -m 0 -w 40 -roll -rra 1 -iter 16```	

Samples : 
```martine -i samples/Batman-Neal-Adams.jpg -m 0 -f```

input ![samples/Batman-Neal-Adams.jpg](samples/Batman-Neal-Adams.jpg)      
 
 will resize the image and save it 

 ![resized](samples/batman_mode0_resized.png)

 after downgrade the colors palette : 

 ![downgrade colors](samples/batman_mode0_down.png)

 results on a CPC emulator : 

 ![result](samples/overscan-batman.png)


Samples roll usage : 

```martine -i samples/rotate.png -m 0 -w 16 -h 16 -roll -rra 1 -iter 16```

input ![samples/rotate.png](samples/rotate.png)

sames phasis reduce size and downgrade colors palette to CPC palette. 

after rotate the first pixels' column in 16 differents images : 

 ![0rotate.png](samples/0rotate.png)
 ![1rotate.png](samples/1rotate.png)
 ![2rotate.png](samples/2rotate.png)
 ![3rotate.png](samples/3rotate.png)
 ![4rotate.png](samples/4rotate.png)
 ![5rotate.png](samples/5rotate.png)
 ![6rotate.png](samples/6rotate.png)
 ![7rotate.png](samples/7rotate.png)
 ![8rotate.png](samples/8rotate.png)
 ![9rotate.png](samples/9rotate.png)
 ![10rotate.png](samples/10rotate.png)
 ![11rotate.png](samples/11rotate.png)
 ![12rotate.png](samples/12rotate.png)
 ![13rotate.png](samples/13rotate.png)
 ![14rotate.png](samples/14rotate.png)
 ![15rotate.png](samples/15rotate.png)
