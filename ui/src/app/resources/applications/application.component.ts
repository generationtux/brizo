import { Environment } from '../environments/environment.component'
import { Pod } from '../pod.component'

export class Application {
  constructor(
    public id:   number,
    public ID:   number,
    public uuid: string,
    public name: string,
    public environments: Environment[],
    public pods: Array<Pod>
  ){}
}
