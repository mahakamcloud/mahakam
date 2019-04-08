FROM centos:7

RUN yum -y install \
    unzip \
    genisoimage-1.1.11 \
    libvirt \
    wget \
    git \
    && yum -y clean all

RUN cd tmp \
    && curl -O https://releases.hashicorp.com/terraform/0.11.8/terraform_0.11.8_linux_amd64.zip \
    && unzip ./terraform_* \
    && mv ./terraform /usr/local/bin/terraform

RUN cd tmp \
    && curl -LO https://github.com/dmacvicar/terraform-provider-libvirt/releases/download/v0.4.4/terraform-provider-libvirt-0.4.4.CentOS_7.x86_64.tar.gz \
    && tar -zxvf ./terraform-provider-* \
    && mkdir -p /root/.terraform.d/plugins \
    && mv ./terraform-provider-libvirt /root/.terraform.d/plugins/ \
    && rm -f ./terraform*

RUN cd tmp \
    && curl -OL https://storage.googleapis.com/golang/go1.11.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.11.5.linux-amd64.tar.gz \
    && rm -rf go1.11.5.linux-amd64.tar.gz \
    && echo "export PATH=$PATH:/usr/local/go/bin" >> /etc/profile

ENV PATH=$PATH:/usr/local/go/bin

ENV GOPATH /go
RUN mkdir -p $GOPATH/{bin,main,pkg,src} && chmod -R 777 $GOPATH
ENV PATH $GOPATH/bin:$PATH
WORKDIR $GOPATH

RUN curl https://raw.githubusercontent.com/golang/dep/v0.5.1/install.sh | sh