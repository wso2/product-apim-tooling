#compdef apimcli

_arguments \
  '1: :->level1' \
  '2: :->level2' \
  '3: :_files'
case $state in
  level1)
    case $words[1] in
      apimcli)
        _arguments '1: :(add-env api export-api export-apis export-app help import-api import-app list remove-env reset-user set version)'
      ;;
      *)
        _arguments '*: :_files'
      ;;
    esac
  ;;
  level2)
    case $words[2] in
      api)
        _arguments '2: :(help lifecycle)'
      ;;
      list)
        _arguments '2: :(apis apps envs help)'
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
