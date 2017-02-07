import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { RouterModule, Routes } from '@angular/router';
import { SharedModule } from '../shared.module';

import { ApplicationCreationComponent } from './creation/application-creation.component';
import { ApplicationDetailsComponent } from './details/application-details.component';
import { ApplicationListingComponent } from './listing/application-listing.component';
import { AuthGuard } from '../auth/auth.guard';
import { EnvironmentService } from '../environments/environment.service';

const applicationRoutes: Routes = [
  { path: 'applications/:uuid', component: ApplicationDetailsComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule.forChild(applicationRoutes),
    SharedModule,
  ],
  declarations: [
    ApplicationCreationComponent,
    ApplicationDetailsComponent,
    ApplicationListingComponent,
  ],
  providers: [ EnvironmentService ],
  exports: [
    ApplicationCreationComponent,
    ApplicationDetailsComponent,
    ApplicationListingComponent,
    RouterModule,
  ],
})

export class ApplicationModule {};
