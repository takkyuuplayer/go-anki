augroup QuickRunGo
  autocmd!

  let g:quickrun_config['go'] = {}
  let g:quickrun_config['go']['command'] = 'make'
  let g:quickrun_config['go']['cmdopt'] = 'run'
  let g:quickrun_config['go']['exec'] = '%c %o SRC=@%'
  let g:quickrun_config['go']['hook/cd/enable'] = 1
  let g:quickrun_config['go']['hook/cd/directory'] = '$PWD'

  autocmd BufWinEnter,BufNewFile *_test.go set filetype=go.test

  let g:quickrun_config['go.test'] = {}
  let g:quickrun_config['go.test']['command'] = 'make'
  let g:quickrun_config['go.test']['cmdopt'] = 'run-test'
  let g:quickrun_config['go.test']['exec'] = '%c %o SRC=@%'
  let g:quickrun_config['go.test']['hook/cd/enable'] = 1
  let g:quickrun_config['go.test']['hook/cd/directory'] = '$PWD'

augroup END
