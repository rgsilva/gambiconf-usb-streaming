# GambiConf 2024

Este é o código referente à talk "Bluetooth pra quê? Fazendo streaming pela porta USB" apresentada na GambiConf 2024.

## Compilando

Primeiramente você precisa ter o golang instalado na sua máquina. O código foi feito para go 1.23, mas provavelmente versões próximas irão funcionar sem problema. Para compilar todos os demos, basta rodar o make:

```sh
make
```

## Clock FS

A primeira demo consiste de um sistema de arquivos que contém apenas um arquivo dentro dele chamado `clock`. Para executar:

```sh
mkdir -p ./mnt
bin/clockfs ./mnt
```

## Corrupt FS

O segundo demo consiste de um pequeno sistema de arquivos virtual que usa uma pasta real por trás e gera um FS contendo os mesmos arquivos. Porém, ao tentar ler esses outros arquivos, o FS corrompe eles em tempo real e de forma totalmente aleatória. Para executar:

```sh
mkdir -p ./mnt
bin/corruptfs files/demo2 ./mnt
```

## Stream FS

O terceiro e último demo é todo o objetivo da talk. Esse FS, uma vez montado, contém apenas um único arquivo, intercepta leituras dentro de uma imagem de disco. Para executar:

```sh
mkdir -p ./mnt
bin/streamfs 5000 files/demo3/imagem.raw ./mnt
```

⚠️ O demo hoje só funciona corretamente com a porta 5000 e tem alguns offsets parcialmente hard-coded (senão ia levar uns 10min na GambiConf para rodar). Talvez você precise ajustar isto caso mude os arquivos. Entretanto, os arquivos de exemplo são fornecidos neste repositório.

Para enviar dados para o FS, você pode um reference stream ou qualquer outro streaming MP3 (ou até mesmo arquivos). Um exemplo é o abaixo:

```sh
curl -L "https://streams.radiomast.io/ref-128k-mp3-stereo" > /dev/tcp/127.0.0.1/5000
```

Não se esqueça de montar o disco também usando o USB Gadget Mass Storage:

```sh
sudo modprobe g_mass_storage file=$(pwd)/mnt/imagem.raw
```

Para testar local, sem utilizar o USB Gadget, também é possível montar a imagem do disco localmente. Isto pode ser feito utilizando funcionalidades normais do Linux, mas forneci também os arquivos `mount.sh` e `umount.sh` junto aos arquivos do demo 3 que auxiliam neste processo.

## Parte elétrica

Nem todo rádio dá conta da Raspberry Pi 4 fazendo boot (2 a 3A). Para isto, na talk eu utilizei um cabo USB-A -> USB-C modificado, onde o pino de 5V foi cortado. Conectado ao ground ("negativo") e ao 5V ("positivo") do lado do USB-C (o lado que vai para a Raspberry Pi), eu conectei uma trigger board de Power Delivery, que pode ser facilmente encontrada no AliExpress (pesquise por "trigger board pd", por exemplo). Essa trigger deve ser configurada para 5V, que é a voltagem esperada pela Pi. O fio dos 5V que vai para a porta USB A **não** deve ser conectado (ou vai possivelmente queimar teu rádio).

Feito isto, a RPi recebia energia por um powerbank, porém os dados ainda passavam normalmente entre ela e o rádio.

## Aviso de Isenção

Este código é fornecido "como está" (as is), sem garantias de qualquer tipo, expressas ou implícitas. O uso é de sua inteira
responsabilidade. Não me responsabilizo por quaisquer danos diretos, indiretos, incidentais ou consequenciais,
decorrentes do uso deste código. Ah, e eu recomendo ter um extintor por perto. Vai que.
