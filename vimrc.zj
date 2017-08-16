
"tab 键等于4个空格的长度
set ts=4

"tab 变成空格
"set expandtab

"自动缩进
set smartindent

"当使用 et 将 Tab 替换为空格之后， 按下一个 Tab 键就能插入 4 个空格，
"但要想删除这 4 个空格， 就得按 4 下 Backspace， 很不方便。 设置 smarttab
"之后， 就可以只按一下 Backspace 就删除 4 个空格了。
set smarttab

"这个是用于程序中自动缩进所使用的空白长度指示的。
set shiftwidth=4

"搜索时高亮显示被找到的文本
set hls

"打开 C/C++ 风格的自动缩进
set cin

"打开关键字上色
syntax on

"设置快速搜索, 输入字符串就显示匹配点
set incsearch


"状态栏
set statusline=[%F]=[Line:%l,Column:%c,Lines:%L][%p%%]
set laststatus=2

let g:pep8_map='whatever'

set backspace=2

"每行最多81个字符
"highlight OverLength ctermbg=red ctermfg=white guibg=#592929
"match OverLength /\%81v.\+/



"记录vim 上一次修改位置
autocmd BufReadPost *
\ if line("'\"")>0&&line("'\"")<=line("$") |
\   exe "normal g'\"" |
\ endif

"  vim 退出后，保留之前看过的内容
set t_ti= t_te=




set nocompatible              " be iMproved, required
filetype off                  " required

" set the runtime path to include Vundle and initialize
set rtp+=~/.vim/bundle/Vundle.vim
call vundle#begin()

" let Vundle manage Vundle, required
Plugin 'gmarik/Vundle.vim'
Plugin 'fatih/vim-go'
Plugin 'Valloric/YouCompleteMe'
Bundle 'dgryski/vim-godef'
Plugin 'nsf/gocode', {'rtp': 'vim/'}

" All of your Plugins must be added before the following line
call vundle#end()            " required
filetype plugin indent on    " required


"vim-godef   reuse the current window, 
"https://github.com/dgryski/vim-godef
" 快捷键 gd  golang 代码跳转
let g:godef_split=0



"autocmd BufWritePre *.go fmt
let g:go_fmt_command = "goimports"
