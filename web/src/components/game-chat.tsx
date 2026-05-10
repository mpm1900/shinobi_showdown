import { InputGroup } from 'node_modules/@base-ui/react/esm/autocomplete/index.parts'
import {
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from './ui/input-group'
import { ArrowRight } from 'lucide-react'

function GameChat() {
  return (
    <div className="flex flex-col">
      <div className="flex-1"></div>
      <div>
        <InputGroup>
          <InputGroupInput placeholder="Send a message" />
          <InputGroupAddon align="inline-end">
            <InputGroupButton variant="default">
              <ArrowRight />
            </InputGroupButton>
          </InputGroupAddon>
        </InputGroup>
      </div>
    </div>
  )
}

export { GameChat }
