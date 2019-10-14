#compdef apimctl

_arguments \
  '1: :->level1' \
  '2: :->level2' \
  '3: :_files'
case $state in
  level1)
    case $words[1] in
      apimctl)
        _arguments '1: :(add add-env export-api export-apis export-app get-keys help import-api import-app init list login logout remove-env set update version)'
      ;;
      *)
        _arguments '*: :_files'
      ;;
    esac
  ;;
  level2)
    case $words[2] in
      add)
        _arguments '2: :(api help)'
      ;;
      list)
        _arguments '2: :(apis apps envs help)'
      ;;
      update)
        _arguments '2: :(api help)'
      ;;
      *)
        _arguments '*: :_files'
      ;;
    esac
  ;;
  *)
    _arguments '*: :_files'
  ;;
esac
