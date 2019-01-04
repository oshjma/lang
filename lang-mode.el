(defconst lang-keywords
  (list
   "let"
   "if"
   "else"
   "for"
   "in"
   "continue"
   "break"
   "return"))

(defconst lang-basictypes
  (list
   "int"
   "bool"
   "string"))

(defconst lang-builtins
  (list
   "true"
   "false"
   "puts"
   "printf"))

(defconst lang-font-lock-keywords
  `(;; Keywords
    (,(rx symbol-start
          (eval `(or ,@lang-keywords))
          symbol-end)
     0 font-lock-keyword-face)

    ;; Basic types
    (,(rx symbol-start
          (eval `(or ,@lang-basictypes))
          symbol-end)
     0 font-lock-type-face)

    ;; Builtins
    (,(rx symbol-start
          (eval `(or ,@lang-builtins))
          symbol-end)
     0 font-lock-builtin-face)

    ;; Let statement
    ;;(,(rx symbol-start
    ;;      "let"
    ;;      (1+ space)
    ;;      (group (1+ (or alnum "_")))
    ;;      symbol-end)
    ;; 1 font-lock-variable-name-face)

    ;; Function parameter
    ;;(,(rx symbol-start
    ;;      (group (1+ (or alnum "_")))
    ;;      (0+ space)
    ;;      ":")
    ;; 1 font-lock-variable-name-face)

    ;; Function call
    (,(rx symbol-start
          (group (1+ (or alnum "_")))
          (0+ space)
          "(")
     1 font-lock-function-name-face)))

(define-derived-mode lang-mode prog-mode "lang"
  "Major mode for editing lang files."
  (setq-local font-lock-defaults '(lang-font-lock-keywords)))

(add-to-list 'auto-mode-alist '("\\.lg\\'" . lang-mode))

(provide 'lang-mode)
