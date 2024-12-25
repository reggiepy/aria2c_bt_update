# compile for version
make
if [ $? -ne 0 ]; then
    echo "make error"
    exit 1
fi

version=`./github.com/reggiepy/aria2c_bt_updater -v`
echo "build version: $version"

# cross_compiles
make -f ./Makefile.cross-compiles

rm -rf ./release/packages
mkdir -p ./release/packages

os_all='linux windows darwin freebsd'
arch_all='386 amd64 arm arm64 mips64 mips64le mips mipsle riscv64'

cd ./release

for os in $os_all; do
    for arch in $arch_all; do
        frp_dir_name="frp_${version}_${os}_${arch}"
        frp_path="./packages/${version}_${os}_${arch}"

        if [ "x${os}" = x"windows" ]; then
            if [ ! -f "./${os}_${arch}.exe" ]; then
                continue
            fi
            if [ ! -f "./${os}_${arch}.exe" ]; then
                continue
            fi
            mkdir ${frp_path}
            mv ./frpc_${os}_${arch}.exe ${frp_path}/frpc.exe
            mv ./frps_${os}_${arch}.exe ${frp_path}/frps.exe
        else
            if [ ! -f "./${os}_${arch}" ]; then
                continue
            fi
            if [ ! -f "./${os}_${arch}" ]; then
                continue
            fi
            mkdir ${frp_path}
            mv ./${os}_${arch} ${frp_path}/frpc
            mv ./${os}_${arch} ${frp_path}/frps
        fi
        cp ../LICENSE ${frp_path}
        cp -rf ../conf/* ${frp_path}

        # packages
        cd ./packages
        if [ "x${os}" = x"windows" ]; then
            zip -rq ${frp_dir_name}.zip ${frp_dir_name}
        else
            tar -zcf ${frp_dir_name}.tar.gz ${frp_dir_name}
        fi
        cd ..
        rm -rf ${frp_path}
    done
done

cd -