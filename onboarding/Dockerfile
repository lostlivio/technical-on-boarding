FROM golang:1.8

ENV PACKAGE_PATH "/go/src/github.com/samsung-cnct/technical-on-boarding/onboarding"


RUN apt-get -qq update && apt-get install -y -q build-essential


WORKDIR "${PACKAGE_PATH}"
COPY "./" "${PACKAGE_PATH}"

RUN make all
RUN mkdir /workload && \
    cp -v ${PACKAGE_PATH}/prepare_workload ${PACKAGE_PATH}/onboarding-issues.yaml /workload/

VOLUME ["/go/"]
EXPOSE 7000


ENTRYPOINT [ \
    "/workload/prepare_workload", \
    "/workload/onboarding-issues.yaml" \
]

