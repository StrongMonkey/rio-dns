package rio

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/rewrite"
	"strings"

	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

type Rio struct{
	Next             plugin.Handler
}

func (r Rio) ServeDNS(ctx context.Context, w dns.ResponseWriter, req *dns.Msg) (int, error) {
	wr := rewrite.NewResponseReverter(w, req)
	state := request.Request{W: w, Req: req}
	fqdn := strings.ToLower(state.QName())
	if strings.HasSuffix(fqdn, "rio.local.") && strings.Count(fqdn, ".") == 5 {
		parts := strings.Split(fqdn, ".")
		name, stack, project := parts[0], parts[1], parts[2]
		answerFqdn := fmt.Sprintf("%s-%s.rio-cloud.svc.cluster.local.", name, stackNamespaceOnlyHash(project, stack))
		req.Question[0].Name = answerFqdn
	}
	return plugin.NextOrFailure(r.Name(), r.Next, ctx, wr, req)
}

func stackNamespaceOnlyHash(projectName, stackName string) string {
	parts := strings.Split(stackName, ":")
	stackName = parts[len(parts)-1]

	id := fmt.Sprintf("%s:%s", projectName, stackName)
	h := sha256.New()
	h.Write([]byte(id))
	hash := hex.EncodeToString(h.Sum(nil))
	return string(hash)[:8]
}

// Name implements the Handler interface.
func (r Rio) Name() string { return "rio" }
