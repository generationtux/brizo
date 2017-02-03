import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AuthGuard } from '../auth/auth.guard';
import { Environment } from './environments/environment.component';
import { MastheadModule } from '../masthead/masthead.module';
import { EnvironmentService } from './environments/environment.service';
import { Application } from './applications/application.component';
import { ApplicationComponent } from './applications/application-details.component';
import { ApplicationCreationComponent } from './applications/application-creation.component';
import { ApplicationListingComponent } from './applications/application-listing.component';
import { EnvironmentDetailsComponent } from './environments/environment-details.component';

@NgModule({
  imports: [
    RouterModule.forChild([
      { path: 'applications/:uuid', component: ApplicationComponent, canActivate: [AuthGuard] },
      { path: 'environments/:uuid', component: EnvironmentDetailsComponent, canActivate: [AuthGuard] },
    ]),
    MastheadModule,
  ],
  declarations: [
    ApplicationComponent,
    ApplicationCreationComponent,
    ApplicationListingComponent,
    EnvironmentDetailsComponent,
  ],
  exports:      [
    ApplicationComponent,
    Application,
    Environment,
    EnvironmentDetailsComponent,
  ],
  providers:    [ EnvironmentService ],
})
export class ResourcesModule { }
