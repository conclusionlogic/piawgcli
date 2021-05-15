BIN_FILE_NAME_PREFIX=$1
PROJECT_DIR=$2
#PLATFORMS=$(go tool dist list)
PLATFORMS="linux/amd64
freebsd/amd64
windows/amd64"
for PLATFORM in $PLATFORMS; do
        GOOS=${PLATFORM%/*}
        GOARCH=${PLATFORM#*/}
        FILEPATH="$PROJECT_DIR/artifacts"
        mkdir -p $FILEPATH
        BIN_FILE_NAME="$FILEPATH/${BIN_FILE_NAME_PREFIX}"
        if [[ "${GOOS}" == "windows" ]]; then BIN_FILE_NAME="${BIN_FILE_NAME}.exe"; fi
        CMD="GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${BIN_FILE_NAME}"
        echo "${CMD}"
        eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
        gzip -c $BIN_FILE_NAME > ${GOOS}-${GOARCH}-${BIN_FILE_NAME}.gz
        rm -f $BIN_FILE_NAME
done