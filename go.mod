module github.com/RPDevJesco/retroanalysis-core

go 1.24

replace github.com/RPDevJesco/retroanalysis-drivers => ../retroanalysis-drivers

require (
    github.com/RPDevJesco/retroanalysis-drivers v0.0.0-00010101000000-000000000000
    cuelang.org/go v0.6.0
    github.com/spf13/viper v1.18.2
)