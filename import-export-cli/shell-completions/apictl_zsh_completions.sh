#compdef apictl

_arguments \
  '1: :->level1' \
  '2: :->level2' \
  '3: :_files'
case $state in
  level1)
    case $words[1] in
      apictl)
        _arguments '1: :(add add-env change export-api export-apis export-app get-keys help import-api import-app init install list login logout remove-env set uninstall update version)'
      ;;
      *)
        _arguments '*: :_files'
      ;;
    esac
  ;;
  level2)
    case $words[2] in
      change)
        _arguments '2: :(help registry)'
      ;;
      install)
        _arguments '2: :(api-operator help)'
      ;;
      list)
        _arguments '2: :(apis apps envs help)'
      ;;
      uninstall)
        _arguments '2: :(api-operator help)'
      ;;
      update)
        _arguments '2: :(api help)'
      ;;
      add)
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
