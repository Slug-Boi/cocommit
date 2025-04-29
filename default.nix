{ lib, buildGoModule, fetchFromGitHub }: 

buildGoModule rec {
  pname = "cocommit";
  version = "1.3.0";

  src = fetchFromGitHub {
    owner = "Slug-Boi";
    repo = pname;
    rev = "v${version}";
    sha256 = "sha256-mSu9IW14y4vgvV3/N4EG9oMvB5eTfcneF03kMmHMXIU=";
  };

  vendorHash = "sha256-GcRGae42KiqMhxc2Q7Ct+uJ4Wg2odUEwWyXffamOjWY=";

  buildPhase = ''
    make build
  '';
  
  doCheck = false;

  installPhase = ''
    mkdir -p $out/bin
    cp "src/${pname}" "$out/bin/${pname}"
    chmod +x $out/bin/${pname}
  '';

  makefile = "makefile";

  meta = with lib; {
    description = "Cocommit is a CLI that makes it easier to co-author users on git commits";
    homepage = "https://github.com/Slug-Boi/${pname}";
    license = licenses.mit;
    platforms = platforms.unix;
  };
}
