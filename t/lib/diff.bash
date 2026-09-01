function @diff() {

        # $1/$2 contain the name of the variable we want to reference from the
        # parent scope

        local -n left=$1
        local -n right=$2

        if type delta &> /dev/null ; then
                #DIFF_COMMAND="delta -w 72  --detect-dark-light always --side-by-side"
                DIFF_COMMAND="delta --detect-dark-light always"
        else
                DIFF_COMMAND="diff -u -w"
        fi

        $DIFF_COMMAND \
                <(printf "%s\n" "${left[@]}") \
                <(printf "%s\n" "${right[@]}")
}

# function wrap_and_fmt_expected_output() {
#         jq --sort-keys '.data |= sort_by(.loginname)' <<<"{ \"data\": [ $@ ] }"
# }

# function fmt_result_output() {
#         jq -r --sort-keys '.data |= sort_by(.loginname)'
# }
