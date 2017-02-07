import { Environment } from '../environments/environment.model';
import { Pod } from '../pod.model';

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
