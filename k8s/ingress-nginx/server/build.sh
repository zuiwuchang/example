
  
#!/usr/bin/env bash

set -e

BashDir=$(cd "$(dirname $BASH_SOURCE)" && pwd)
eval $(cat "$BashDir/script/conf.sh")
if [[ "$Command" == "" ]];then
    Command="$0"
fi

function help(){
    echo "build script"
    echo
    echo "Usage:"
    echo "  $0 [flags]"
    echo "  $0 [command]"
    echo
    echo "Available Commands:"
    echo "  help              help for $0"
    echo "  clear             clear output"
    echo "  go                go build helper"
    echo "  grpc              grpc protoc helper"
    echo "  pack              pack release"
    echo "  run               run project"
    echo "  docker            docker build helper"
    echo
    echo "Flags:"
    echo "  -h, --help          help for $0"
}

case "$1" in
    help|-h|--help)
        help
    ;;
    clear)
        shift
        export Command="$0 clear"
        "$BashDir/script/clear.sh" "$@"
    ;;
    pack)
        shift
        export Command="$0 pack"
        "$BashDir/script/pack.sh" "$@"
    ;;
    go)
        shift
        export Command="$0 go"
        "$BashDir/script/go.sh" "$@"
    ;;
    grpc)
        shift
        export Command="$0 grpc"
        "$BashDir/script/grpc.sh" "$@"
    ;;
    run)
        shift
        export Command="$0 run"
        "$BashDir/script/run.sh" "$@"
    ;;
    docker)
        shift
        export Command="$0 docker"
        "$BashDir/script/docker.sh" "$@"
    ;;
    *)
        if [[ "$1" == "" ]];then
            help
        elif [[ "$1" == -* ]];then
            echo Error: unknown flag "$1" for "$0"
            echo "Run '$0 --help' for usage."
        else
            echo Error: unknown command "$1" for "$0"
            echo "Run '$0 --help' for usage."
        fi        
        exit 1
    ;;
esac