import React from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router-dom';
import { Link as UswdsLink } from '@trussworks/react-uswds';
import { DateTime } from 'luxon';

import PageHeading from 'components/PageHeading';
import PDFExport from 'components/PDFExport';
import { AnythingWrongSurvey } from 'components/Survey';
import SystemIntakeReview from 'components/SystemIntakeReview';
import { SystemIntakeForm } from 'types/systemIntake';

type IntakeReviewProps = {
  systemIntake: SystemIntakeForm;
  now: DateTime;
};

const IntakeReview = ({ systemIntake, now }: IntakeReviewProps) => {
  const { t } = useTranslation('governanceReviewTeam');
  const filename = `System intake for ${systemIntake.requestName}.pdf`;

  return (
    <div>
      <PageHeading className="margin-top-0">{t('general:intake')}</PageHeading>
      <PDFExport
        title="System Intake"
        filename={filename}
        label="Download System Intake as PDF"
      >
        <SystemIntakeReview systemIntake={systemIntake} now={now} />
      </PDFExport>
      <UswdsLink
        className="usa-button margin-top-5"
        variant="unstyled"
        to={`/governance-review-team/${systemIntake.id}/actions`}
        asCustom={Link}
      >
        Take an action
      </UswdsLink>
      <AnythingWrongSurvey />
    </div>
  );
};

export default IntakeReview;
