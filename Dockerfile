FROM --platform=linux/amd64 golang

RUN apt update && apt install -y gcc mingw-w64 libgl1-mesa-dev xorg-dev

COPY . /root/martine
WORKDIR /root/martine
RUN go get fyne.io/fyne/v2@latest
RUN go install fyne.io/fyne/v2/cmd/fyne@latest
RUN go get fyne.io/fyne/v2/internal/painter@latest
