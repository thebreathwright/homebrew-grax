package renderers

import (
	"strings"
	"unicode"

	"github.com/ollama/ollama/api"
)

const graxDefaultSystem = "You are the Brullama doer. Your operating roles are reader, listener, and do. Bind role tokens to an explicit recipient and subject with evidence; a floating role token is RECIPIENT_NOT_ESTABLISHED. Reader preserves source identity and authority boundaries. Listener begins from what the runtime cannot yet know: inventory the actual observation vector O*, then construct valid observationally identical realities before naming a distinction. Do not add fields, validators, receipts, or production patches until a collision is proven to require different lawful treatment. Constructor-supplied outcome labels are STIPULATED, not evidence. Require a provenance-bound consequence basis and outcome-swap check before COLLISION_CONFIRMED. Derive consequences independently for the admitted observation class. Execute C only when its consequence set is a singleton; otherwise withhold C as NOT_IDENTIFIABLE or INDEPENDENCE_UNRESOLVED. Choose MINIMUM_MISSING_OBSERVATION by consequence heterogeneity reduced per declared cost. Do not invent a classifier. Test false admission and false withholding symmetrically. Candidate-set digests bind only the DECLARED_CANDIDATES_ONLY evaluation scope; do not treat an evaluated set as complete without separate evidence. Report operation success separately from WorkDelta; unchanged uncertainty means WorkDelta=0. Reader preserves requested/bound/effective parameters, source meaning, observed transformation, and unresolved meaning. A crystal may guide observation but cannot donate authority. Representation does not establish meaning; output presence does not establish transfer; a hash does not establish truth. Do runs only explicitly authorized local operations and freezes receipts. Produce concise operational audit trails with answer, evidence_status, boundary, next_safe_action, and work_delta. Example answer shape: answer=..., evidence_status=NOT_ESTABLISHED|SUPPORTED|UNSUPPORTED, boundary=..., next_safe_action=..., work_delta=...."

const graxQueueContract = "Queued deictic or low-information instructions such as 'go make that' and 'run it' must preserve submission-time referent and goal identity separately from their text. If the target is absent, return TARGET_UNRESOLVED; repetition does not create referential completeness. Do not silently bind a target from contextual proximity, interpretation, or a predecessor's result. Resolve an operation only under an explicit continuation capability with a permitted operation class, consequence envelope, revalidation rule, transition budget, stop condition, and expiry. If multiple candidates remain, preserve TARGET_UNRESOLVED rather than choosing by recency or convenience. Distinguish FROZEN_REFERENT from BOUNDED_CONTINUATION: revalidation may confirm, reject, or leave a binding unresolved, but may not replace a frozen referent with a newly salient object. A bounded continuation may advance only through its authorized successor relation. Zero work, proof, and consequence deltas are PROGRESS_STUTTERING, not advancement. Continuation preserves the goal, not the prior operation or strategy: reselect from current state only within that capability. Duplicate pending continuation capabilities are idempotent and do not amplify mutation authority; intervening state changes may revoke them."

type GraxRenderer struct{}

func (r *GraxRenderer) LeadingBOS() string {
	return lagunaBOS
}

func (r *GraxRenderer) Render(messages []api.Message, tools []api.Tool, think *api.ThinkValue) (string, error) {
	var sb strings.Builder
	sb.WriteString(lagunaBOS)

	thinkingEnabled := think != nil && think.Bool()

	systemMessage := graxDefaultSystem + "\n\n" + graxQueueContract
	firstMessageIsSystem := len(messages) > 0 && messages[0].Role == "system"
	if firstMessageIsSystem {
		systemMessage = messages[0].Content
	}

	if strings.TrimSpace(systemMessage) != "" || len(tools) > 0 {
		sb.WriteString("<system>\n")
		if strings.TrimSpace(systemMessage) != "" {
			sb.WriteByte('\n')
			sb.WriteString(strings.TrimRightFunc(systemMessage, unicode.IsSpace))
		}
		if len(tools) > 0 {
			sb.WriteString("\n\n### Tools\n\n")
			sb.WriteString("You may call functions to assist with the user query.\n")
			sb.WriteString("All available function signatures are listed below:\n")
			sb.WriteString("<available_tools>\n")
			for _, tool := range tools {
				if b, err := marshalWithSpaces(tool); err == nil {
					sb.Write(b)
					sb.WriteByte('\n')
				}
			}
			sb.WriteString("</available_tools>\n\n")
			if thinkingEnabled {
				sb.WriteString("Wrap your thinking in '<think>', '</think>' tags, then emit the most bounded valid tool call or answer. Keep the output operational and claim-bounded. Use the reader/listener/do boundaries to avoid authority inflation and false passes. Prefer answer, evidence_status, boundary, and next_safe_action when you answer directly.\n")
				sb.WriteString("<think> your thoughts here </think>\n")
			} else {
				sb.WriteString("For each function call, return an unescaped XML-like object with function name and arguments within '<tool_call>' and '</tool_call>' tags, like here:\n")
			}
			sb.WriteString("<tool_call>function-name\n<arg_key>argument-key</arg_key>\n<arg_value>value-of-argument-key</arg_value>\n</tool_call>")
		}
		sb.WriteString("\n</system>\n")
	}

	for i, message := range messages {
		if i == 0 && firstMessageIsSystem {
			continue
		}
		content := message.Content
		switch message.Role {
		case "user":
			sb.WriteString("<user>\n")
			sb.WriteString(content)
			if think != nil && !think.Bool() {
				sb.WriteString(" /no_think")
			}
			sb.WriteString("\n</user>\n")
		case "assistant":
			content, reasoning := lagunaV2AssistantContent(message.Content, message.Thinking)
			lastMessage := i == len(messages)-1
			prefill := lastMessage && (strings.TrimSpace(content) != "" || strings.TrimSpace(reasoning) != "" || len(message.ToolCalls) > 0)

			sb.WriteString("<assistant>\n")
			if reasoning := strings.TrimSpace(reasoning); reasoning != "" {
				sb.WriteString("<think>\n")
				sb.WriteString(reasoning)
				sb.WriteString("\n</think>\n")
			} else {
				sb.WriteString("</think>\n")
			}

			if strings.TrimSpace(content) != "" {
				sb.WriteString(strings.TrimSpace(content))
				sb.WriteByte('\n')
			}

			for _, toolCall := range message.ToolCalls {
				sb.WriteString("<tool_call>")
				sb.WriteString(toolCall.Function.Name)
				sb.WriteByte('\n')
				for name, value := range toolCall.Function.Arguments.All() {
					sb.WriteString("<arg_key>")
					sb.WriteString(name)
					sb.WriteString("</arg_key>\n")
					sb.WriteString("<arg_value>")
					sb.WriteString(formatLagunaToolCallArgument(value))
					sb.WriteString("</arg_value>\n")
				}
				sb.WriteString("</tool_call>\n")
			}

			if !prefill {
				sb.WriteString("</assistant>\n")
			}
		case "tool":
			sb.WriteString("<tool_response>\n")
			sb.WriteString(content)
			sb.WriteString("\n</tool_response>\n")
		case "system":
			sb.WriteString("<system>\n")
			sb.WriteString(content)
			sb.WriteString("\n</system>\n")
		}
	}

	if len(messages) == 0 || messages[len(messages)-1].Role != "assistant" {
		sb.WriteString("<assistant>\n")
		if thinkingEnabled {
			sb.WriteString(lagunaThoughtOpen)
		} else {
			sb.WriteString(lagunaThoughtOpen)
			sb.WriteString("\n\n")
			sb.WriteString(lagunaThoughtClose)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}
