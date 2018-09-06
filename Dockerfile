####################################################
# GOLANG BUILDER
####################################################
FROM golang:1.11 as go_builder

COPY . /go/src/github.com/malice-plugins/comodo
WORKDIR /go/src/github.com/malice-plugins/comodo
RUN go get -u github.com/golang/dep/cmd/dep && dep ensure
RUN go build -ldflags "-s -w -X main.Version=v$(cat VERSION) -X main.BuildTime=$(date -u +%Y%m%d)" -o /bin/avscan

####################################################
# PLUGIN BUILDER
####################################################
FROM ubuntu:precise

LABEL maintainer "https://github.com/blacktop"

LABEL malice.plugin.repository = "https://github.com/malice-plugins/comodo.git"
LABEL malice.plugin.category="av"
LABEL malice.plugin.mime="*"
LABEL malice.plugin.docker.engine="*"

# Create a malice user and group first so the IDs get set the same way, even as
# the rest of this may change over time.
RUN groupadd -r malice \
  && useradd --no-log-init -r -g malice malice \
  && mkdir /malware \
  && chown -R malice:malice /malware

RUN buildDeps='ca-certificates \
  build-essential \
  gdebi-core \
  libssl-dev \
  mercurial \
  git-core \
  wget' \
  && apt-get update -qq \
  && apt-get install -yq $buildDeps \
  && echo "===> Install Comodo..." \
  && cd /tmp \
  && wget http://download.comodo.com/cis/download/installs/linux/cav-linux_x64.deb \
  && DEBIAN_FRONTEND=noninteractive gdebi -n cav-linux_x64.deb \
  && DEBIAN_FRONTEND=noninteractive /opt/COMODO/post_setup.sh \
  && echo "===> Clean up unnecessary files..." \
  && apt-get purge -y --auto-remove $buildDeps && apt-get clean \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /root/.gnupg

# Ensure ca-certificates is installed for elasticsearch to use https
RUN apt-get update -qq && apt-get install -yq --no-install-recommends ca-certificates \
  && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Update Comodo definitions
RUN mkdir -p /opt/malice
ADD http://download.comodo.com/av/updates58/sigs/bases/bases.cav /opt/COMODO/scanners/bases.cav

# Add EICAR Test Virus File to malware folder
ADD http://www.eicar.org/download/eicar.com.txt /malware/EICAR

COPY --from=go_builder /bin/avscan /bin/avscan

WORKDIR /malware

ENTRYPOINT ["/bin/avscan"]
CMD ["--help"]

####################################################
####################################################